package instana_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"

	"github.com/stretchr/testify/require"
)

const resourceAlertingChannelEmailDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_email" "example" {
  name = "name {{ITERATOR}}"
  emails = [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name {{ITERATOR}} suffix",
	"kind"   : "EMAIL",
	"emails" : [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelEmailDefinition = "instana_alerting_channel_email.example"

func TestCRUDOfAlertingChannelEmailResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelEmailApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.AlertingChannelsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingChannelEmailServerResponseTemplate, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelEmailDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	emailAddress1 := "EMAIL1"
	emailAddress2 := "EMAIL2"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 0), emailAddress1),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 1), emailAddress2),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 0), emailAddress1),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 1), emailAddress2),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceAlertingChannelEmailDefinition(t *testing.T) {
	resource := NewAlertingChannelEmailResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelEmailFieldEmails)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelEmail(t *testing.T) {
	name := NewAlertingChannelEmailResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_channel_email", name, "Expected resource name to be instana_alerting_channel_email")
}

func TestAlertingChannelEmailResourceShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelEmailResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelEmailShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewAlertingChannelEmailResourceHandle()

	require.Equal(t, 1, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
}

func TestShouldReturnStateOfAlertingChannelEmailUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	emails := []interface{}{"email1", "email2"}
	name := "name"
	fullname := "fullname"
	rawData := make(map[string]interface{})
	rawData[AlertingChannelFieldName] = name
	rawData[AlertingChannelFieldFullName] = fullname
	rawData[AlertingChannelEmailFieldEmails] = emails
	meta := "dummy"
	ctx := context.Background()

	result, err := NewAlertingChannelEmailResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, rawData, result)
}

func TestShouldUpdateResourceStateForAlertingChannelEmail(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   "prefix name suffix",
		Emails: []string{"email1", "email2"},
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "prefix name suffix", resourceData.Get(AlertingChannelFieldFullName))

	emails := resourceData.Get(AlertingChannelEmailFieldEmails).(*schema.Set)
	require.Equal(t, 2, emails.Len())
	require.Contains(t, emails.List(), "email1")
	require.Contains(t, emails.List(), "email2")
}

func TestShouldConvertStateOfAlertingChannelEmailToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	emails := []string{"email1", "email2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelEmailFieldEmails, emails)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Len(t, model.(*restapi.AlertingChannel).Emails, 2)
	require.Contains(t, model.(*restapi.AlertingChannel).Emails, "email1")
	require.Contains(t, model.(*restapi.AlertingChannel).Emails, "email2")
}
