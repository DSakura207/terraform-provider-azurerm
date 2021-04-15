package security

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// SecureScoreControlsClient is the API spec for Microsoft.Security (Azure Security Center) resource provider
type SecureScoreControlsClient struct {
	BaseClient
}

// NewSecureScoreControlsClient creates an instance of the SecureScoreControlsClient client.
func NewSecureScoreControlsClient(subscriptionID string, ascLocation string) SecureScoreControlsClient {
	return NewSecureScoreControlsClientWithBaseURI(DefaultBaseURI, subscriptionID, ascLocation)
}

// NewSecureScoreControlsClientWithBaseURI creates an instance of the SecureScoreControlsClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewSecureScoreControlsClientWithBaseURI(baseURI string, subscriptionID string, ascLocation string) SecureScoreControlsClient {
	return SecureScoreControlsClient{NewWithBaseURI(baseURI, subscriptionID, ascLocation)}
}

// List get all security controls within a scope
// Parameters:
// expand - oData expand. Optional.
func (client SecureScoreControlsClient) List(ctx context.Context, expand ExpandControlsEnum) (result SecureScoreControlListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SecureScoreControlsClient.List")
		defer func() {
			sc := -1
			if result.sscl.Response.Response != nil {
				sc = result.sscl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.SecureScoreControlsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.sscl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "List", resp, "Failure sending request")
		return
	}

	result.sscl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "List", resp, "Failure responding to request")
		return
	}
	if result.sscl.hasNextLink() && result.sscl.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client SecureScoreControlsClient) ListPreparer(ctx context.Context, expand ExpandControlsEnum) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(string(expand)) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/secureScoreControls", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client SecureScoreControlsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client SecureScoreControlsClient) ListResponder(resp *http.Response) (result SecureScoreControlList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client SecureScoreControlsClient) listNextResults(ctx context.Context, lastResults SecureScoreControlList) (result SecureScoreControlList, err error) {
	req, err := lastResults.secureScoreControlListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client SecureScoreControlsClient) ListComplete(ctx context.Context, expand ExpandControlsEnum) (result SecureScoreControlListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SecureScoreControlsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, expand)
	return
}

// ListBySecureScore get all security controls for a specific initiative within a scope
// Parameters:
// secureScoreName - the initiative name. For the ASC Default initiative, use 'ascScore' as in the sample
// request below.
// expand - oData expand. Optional.
func (client SecureScoreControlsClient) ListBySecureScore(ctx context.Context, secureScoreName string, expand ExpandControlsEnum) (result SecureScoreControlListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SecureScoreControlsClient.ListBySecureScore")
		defer func() {
			sc := -1
			if result.sscl.Response.Response != nil {
				sc = result.sscl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.SecureScoreControlsClient", "ListBySecureScore", err.Error())
	}

	result.fn = client.listBySecureScoreNextResults
	req, err := client.ListBySecureScorePreparer(ctx, secureScoreName, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "ListBySecureScore", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListBySecureScoreSender(req)
	if err != nil {
		result.sscl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "ListBySecureScore", resp, "Failure sending request")
		return
	}

	result.sscl, err = client.ListBySecureScoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "ListBySecureScore", resp, "Failure responding to request")
		return
	}
	if result.sscl.hasNextLink() && result.sscl.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListBySecureScorePreparer prepares the ListBySecureScore request.
func (client SecureScoreControlsClient) ListBySecureScorePreparer(ctx context.Context, secureScoreName string, expand ExpandControlsEnum) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"secureScoreName": autorest.Encode("path", secureScoreName),
		"subscriptionId":  autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(string(expand)) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/secureScores/{secureScoreName}/secureScoreControls", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListBySecureScoreSender sends the ListBySecureScore request. The method will close the
// http.Response Body if it receives an error.
func (client SecureScoreControlsClient) ListBySecureScoreSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListBySecureScoreResponder handles the response to the ListBySecureScore request. The method always
// closes the http.Response Body.
func (client SecureScoreControlsClient) ListBySecureScoreResponder(resp *http.Response) (result SecureScoreControlList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listBySecureScoreNextResults retrieves the next set of results, if any.
func (client SecureScoreControlsClient) listBySecureScoreNextResults(ctx context.Context, lastResults SecureScoreControlList) (result SecureScoreControlList, err error) {
	req, err := lastResults.secureScoreControlListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listBySecureScoreNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListBySecureScoreSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listBySecureScoreNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListBySecureScoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SecureScoreControlsClient", "listBySecureScoreNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListBySecureScoreComplete enumerates all values, automatically crossing page boundaries as required.
func (client SecureScoreControlsClient) ListBySecureScoreComplete(ctx context.Context, secureScoreName string, expand ExpandControlsEnum) (result SecureScoreControlListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SecureScoreControlsClient.ListBySecureScore")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListBySecureScore(ctx, secureScoreName, expand)
	return
}
