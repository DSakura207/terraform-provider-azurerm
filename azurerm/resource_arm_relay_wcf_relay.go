package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmWCFRelay() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmWCFRelayCreateUpdate,
		Read:   resourceArmWCFRelayRead,
		Update: resourceArmWCFRelayCreateUpdate,
		Delete: resourceArmWCFRelayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"requires_client_authorization": {
				Type:     schema.TypeBool,
				Default:  true,
				ForceNew: true,
				Optional: true,
			},

			"requires_transport_security": {
				Type:     schema.TypeBool,
				Default:  true,
				ForceNew: true,
				Optional: true,
			},

			"relay_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(relay.HTTP),
					string(relay.NetTCP),
				}, false),
			},

			"user_metadata": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmWCFRelayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Relay.WCFRelaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for WCF Relay creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	relayNamespace := d.Get("relay_namespace_name").(string)
	relayType := d.Get("relay_type").(relay.RelaytypeEnum)
	requireClientAuthroization := d.Get("requires_client_authorization").(bool)
	requiresTransportSecurity := d.Get("requires_transport_security").(bool)
	userMetadata := d.Get("user_metadata").(string)

	parameters := relay.WcfRelay{
		WcfRelayProperties: &relay.WcfRelayProperties{
			RelayType:                   relayType,
			RequiresClientAuthorization: &requireClientAuthroization,
			RequiresTransportSecurity:   &requiresTransportSecurity,
			UserMetadata:                &userMetadata,
		},
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, relayNamespace, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating WCF Relay %q (Namespace %q Resource Group %q): %+v", name, relayNamespace, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, relayNamespace, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for WCF Relay %q (Namespace %q Resource Group %q): %+v", name, relayNamespace, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read WCF Relay %q (Namespace %q Resource group %s) ID", name, relayNamespace, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmWCFRelayRead(d, meta)
}

func resourceArmWCFRelayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Relay.WCFRelaysClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	relayNamespace := id.Path["namespaces"]
	name := id.Path["wcfRelays"]

	resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on WCF Relay %q (Namespace %q Resource Group %q): %s", name, relayNamespace, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("relay_namespace_name", relayNamespace)

	if props := resp.WcfRelayProperties; props != nil {
		d.Set("requires_client_authorization", props.RequiresClientAuthorization)
		d.Set("requires_transport_security", props.RequiresTransportSecurity)
		d.Set("user_metadata", props.UserMetadata)
		d.Set("relay_type", props.RelayType)
	}

	return nil
}

func resourceArmWCFRelayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Relay.WCFRelaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	relayNamespace := id.Path["namespaces"]
	name := id.Path["wcfRelays"]

	log.Printf("[INFO] Waiting for WCF Relay %q (Namespace %q Resource Group %q) to be deleted", name, relayNamespace, resourceGroup)
	rc, err := client.Delete(ctx, resourceGroup, relayNamespace, name)

	if err != nil {
		if response.WasNotFound(rc.Response) {
			return nil
		}

		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    wcfRelayDeleteRefreshFunc(ctx, client, resourceGroup, relayNamespace, name),
		MinTimeout: 15 * time.Second,
	}

	if features.SupportsCustomTimeouts() {
		stateConf.Timeout = d.Timeout(schema.TimeoutDelete)
	} else {
		stateConf.Timeout = 30 * time.Minute
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for WCF Relay %q (Namespace %q Resource Group %q) to be deleted: %s", name, relayNamespace, resourceGroup, err)
	}

	return nil
}

func wcfRelayDeleteRefreshFunc(ctx context.Context, client *relay.WCFRelaysClient, resourceGroupName string, relayNamespace string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, relayNamespace, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}

			return nil, "Error", fmt.Errorf("Error issuing read request in resourceArmWCFRelayDelete to WCF Relay %q (Namespace %q Resource Group %q): %s", name, relayNamespace, resourceGroupName, err)
		}

		return res, "Pending", nil
	}
}
