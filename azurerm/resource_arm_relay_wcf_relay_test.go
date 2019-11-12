package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMRelayWCFRelay_basic(t *testing.T) {
	resourceName := "azurerm_relay_wcf_relay.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayWCFRelayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayWCFRelay_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayWCFRelayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "relay_type"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_transport_security"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRelayWCFRelay_full(t *testing.T) {
	resourceName := "azurerm_relay_wcf_relay.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayWCFRelayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayWCFRelay_full(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayWCFRelayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "relay_type"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_transport_security"),
					resource.TestCheckResourceAttr(resourceName, "user_metadata", "metadatatest"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRelayWCFRelay_update(t *testing.T) {
	resourceName := "azurerm_relay_wcf_relay.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayWCFRelayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayWCFRelay_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayWCFRelayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "relay_type"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_transport_security"),
				),
			},
			{
				Config: testAccAzureRMRelayWCFRelay_update(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relay_type", "NetTcp"),
					resource.TestCheckResourceAttr(resourceName, "requires_client_authorization", "false"),
					resource.TestCheckResourceAttr(resourceName, "requires_transport_security", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_metadata", "metadataupdated"),
				),
			},
		},
	})
}

func TestAccAzureRMRelayWCFRelay_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_relay_wcf_relay.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayWCFRelayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayWCFRelay_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayWCFRelayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "relay_type"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttrSet(resourceName, "requires_transport_security"),
				),
			},
			{
				Config:      testAccAzureRMRelayWCFRelay_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_relay_wcf_relay"),
			},
		},
	})
}

func testAccAzureRMRelayWCFRelay_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_wcf_relay" "test" {
	name                 = "acctestrnhc-%d"
	relay_type           = "Http"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayWCFRelay_full(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_wcf_relay" "test" {
	name                 = "acctestrnhc-%d"
	relay_type           = "Http"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
	user_metadata        = "metadatatest"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayWCFRelay_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_wcf_relay" "test" {
	name                          = "acctestrnhc-%d"
	resource_group_name           = "${azurerm_resource_group.test.name}"
	relay_type                    = "NetTcp"
	relay_namespace_name          = "${azurerm_relay_namespace.test.name}"
	requires_client_authorization = false
	requires_client_authorization = false
	requires_transport_security   = false
	user_metadata                 = "metadataupdated"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayWCFRelay_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_namespace" "import" {
	name                 = "acctestrnhc-%d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
}
`, testAccAzureRMRelayWCFRelay_basic(rInt, location), rInt)
}

func testCheckAzureRMRelayWCFRelayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]

		// Ensure resource group exists in API
		client := testAccProvider.Meta().(*ArmClient).Relay.WCFRelaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on relayWCFRelaysClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: WCF Relay %q in Namespace %q (Resource Group: %q) does not exist", name, relayNamespace, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMRelayWCFRelayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Relay.WCFRelaysClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_relay_wcf_relay" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]

		resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("WCF Relay still exists:\n%#v", resp)
		}
	}

	return nil
}
