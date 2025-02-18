package instana_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_rollup = "40000"
  rule_condition_operator = "="
  rule_condition_value = "1.2"
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_window = 60000
  rule_aggregation = "sum"
  rule_condition_operator = "{{CONDITION_OPERATOR}}"
  rule_condition_value = 1.2
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndMetricPatternDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_window = 60000
  rule_aggregation = "sum"
  rule_condition_operator = "="
  rule_condition_value = 1.2
  rule_metric_pattern_prefix = "prefix"
  rule_metric_pattern_postfix = "postfix"
  rule_metric_pattern_placeholder = "placeholder"
  rule_metric_pattern_operator = "startsWith"
}
`

const (
	testCustomEventSpecificationWithThresholdRuleDefinition = "instana_custom_event_spec_threshold_rule.example"

	customEventSpecificationWithThresholdRuleID             = "custom-system-event-id"
	customEventSpecificationWithThresholdRuleEntityType     = "entity_type"
	customEventSpecificationWithThresholdRuleQuery          = "query"
	customEventSpecificationWithThresholdRuleExpirationTime = 60000
	customEventSpecificationWithThresholdRuleDescription    = "description"
	customEventSpecificationWithThresholdRuleMetricName     = "metric_name"
	customEventSpecificationWithThresholdRuleRollup         = 40000
	customEventSpecificationWithThresholdRuleWindow         = 60000
	customEventSpecificationWithThresholdRuleAggregation    = restapi.AggregationSum
	customEventSpecificationWithThresholdRuleConditionValue = float64(1.2)
)

var CustomEventSpecificationWithThresholdRuleRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithRollupResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName" : "metric_name", "rollup" : 40000, "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldRollup, strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate, "{{CONDITION_OPERATOR}}", "="),
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowAndAlternativeConditionOperatorRepresentationResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate, "{{CONDITION_OPERATOR}}", "=="),
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithMetricPatternResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2, "metricPattern" : { "prefix" : "prefix", "postfix" : "postfix", "placeholder" : "placeholder", "operator" : "startsWith" } }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndMetricPatternDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPrefix, "prefix"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPostfix, "postfix"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPlaceholder, "placeholder"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternOperator, string(restapi.MetricPatternOperatorTypeStartsWith)),
	)
}

const httpServerResponseTemplate = `
{
	"id" : "%s",
	"name" : "prefix name %d suffix",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [ %s ]
}
`

func testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(t *testing.T, terraformDefinition, ruleAsJson string, ruleTestCheckFunctions ...resource.TestCheckFunc) {
	httpServer := createMockHttpServerForResource(restapi.CustomEventSpecificationResourcePath, httpServerResponseTemplate, ruleAsJson)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appendProviderConfig(fmt.Sprintf(terraformDefinition, 0), httpServer.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 0)...),
			},
			testStepImport(testCustomEventSpecificationWithThresholdRuleDefinition),
			{
				Config: appendProviderConfig(fmt.Sprintf(terraformDefinition, 1), httpServer.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 1)...),
			},
			testStepImport(testCustomEventSpecificationWithThresholdRuleDefinition),
		},
	})
}

func createTestCheckFunctions(ruleTestCheckFunctions []resource.TestCheckFunc, iteration int) []resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleDefinition, "id"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldName, formatResourceName(iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldFullName, formatResourceFullName(iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldTriggering, trueAsString),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEnabled, trueAsString),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
	}
	allFunctions := append(defaultCheckFunctions, ruleTestCheckFunctions...)
	return allFunctions
}

func TestCustomEventSpecificationWithThresholdRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithThresholdRuleResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)

	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldWindow)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldRollup)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ThresholdRuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(ThresholdRuleFieldConditionValue)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldMetricPatternPrefix)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldMetricPatternPostfix)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldMetricPatternPlaceholder)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldMetricPatternOperator)
}

func TestCustomEventSpecificationWithThresholdRuleResourceShouldHaveSchemaVersionThree(t *testing.T) {
	require.Equal(t, 3, NewCustomEventSpecificationWithThresholdRuleResourceHandle().MetaData().SchemaVersion)
}

func TestCustomEventSpecificationWithThresholdRuleShouldHaveThreeStateUpgraderForVersionZeroAndOneAndTwo(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	require.Equal(t, 3, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
	require.Equal(t, 2, resourceHandler.StateUpgraders()[2].Version)
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	rawData := make(map[string]interface{})
	rawData[CustomEventSpecificationFieldName] = name
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, name, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateEmptyCustomEventSpecificationWithThresholdRuleStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion2WhenDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["downstream_integration_ids"] = []interface{}{"id1", "id2"}
	rawData["downstream_broadcast_to_all_alerting_configs"] = true
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion2WhenNoDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion3WhenConditionOperatorIsDefinedAndValid(t *testing.T) {
	for _, op := range restapi.SupportedConditionOperators {
		for _, v := range op.TerraformSupportedValues() {
			t.Run(fmt.Sprintf("TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion3WhenConditionOperatorIsDefinedAndValid%s", v), createTestCaseForSuccessfulMigrationOfCustomEventSpecificationWithThresholdRuleToVersion3(op, v))
		}
	}
}

func createTestCaseForSuccessfulMigrationOfCustomEventSpecificationWithThresholdRuleToVersion3(mo restapi.ConditionOperator, value string) func(*testing.T) {
	return func(t *testing.T) {
		rawData := make(map[string]interface{})
		rawData[ThresholdRuleFieldConditionOperator] = value
		meta := "dummy"
		ctx := context.Background()

		result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

		require.Nil(t, err)
		require.Equal(t, mo.InstanaAPIValue(), result[ThresholdRuleFieldConditionOperator])
	}
}

func TestShouldDoNothingWhenMigratingCustomEventSpecificationWithThresholdRuleToVersion3AndNoConditionOperatorIsDefined(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[ThresholdRuleFieldConditionOperator])
}

func TestShouldReturnErrorWhenCustomEventSpecificationWithThresholdRuleCannotBeMigratedToVersion3BecuaseOfUnsupportedConditionOperatorInState(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[ThresholdRuleFieldConditionOperator] = "invalid"
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "not a supported condition operator")
	require.Equal(t, rawData, result)
}

func TestShouldReturnCorrectResourceNameForCustomEventSpecificationWithThresholdRuleResource(t *testing.T) {
	name := NewCustomEventSpecificationWithThresholdRuleResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_custom_event_spec_threshold_rule")
}

func TestShouldUpdateCustomEventSpecificationWithThresholdRuleTerraformStateFromApiObject(t *testing.T) {
	testMappingOfCustomEventSpecificationWithThresholdRuleTerraformDataModelToState(t, func(spec *restapi.CustomEventSpecification) { /* Default testcase without additional fields =< no additional mappings */
	}, func(resourceData *schema.ResourceData) { /* Default testcase without additional fields => no additional requires */
	})
}

func TestShouldUpdateCustomEventSpecificationWithThresholdRuleAndMetricPatternTerraformStateFromApiObject(t *testing.T) {
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := restapi.MetricPatternOperatorTypeStartsWith

	additionalMappings := func(spec *restapi.CustomEventSpecification) {
		metricPattern := restapi.MetricPattern{
			Prefix:      prefix,
			Postfix:     &postfix,
			Placeholder: &placeholder,
			Operator:    operator,
		}
		spec.Rules[0].MetricPattern = &metricPattern
	}

	additionalAsserts := func(resourceData *schema.ResourceData) {
		require.Equal(t, prefix, resourceData.Get(ThresholdRuleFieldMetricPatternPrefix))
		require.Equal(t, postfix, resourceData.Get(ThresholdRuleFieldMetricPatternPostfix))
		require.Equal(t, placeholder, resourceData.Get(ThresholdRuleFieldMetricPatternPlaceholder))
		require.Equal(t, string(operator), resourceData.Get(ThresholdRuleFieldMetricPatternOperator))
	}

	testMappingOfCustomEventSpecificationWithThresholdRuleTerraformDataModelToState(t, additionalMappings, additionalAsserts)
}

func testMappingOfCustomEventSpecificationWithThresholdRuleTerraformDataModelToState(t *testing.T, additionalMappings func(spec *restapi.CustomEventSpecification), additionalAsserts func(resourceData *schema.ResourceData)) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	metricName := customEventSpecificationWithThresholdRuleMetricName
	conditionOperator := restapi.ConditionOperatorEquals.InstanaAPIValue()

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceFullName,
		EntityType:     customEventSpecificationWithThresholdRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          restapi.SeverityWarning.GetAPIRepresentation(),
				MetricName:        &metricName,
				Window:            &window,
				Rollup:            &rollup,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}
	additionalMappings(spec)

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithThresholdRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, resourceFullName, resourceData.Get(CustomEventSpecificationFieldFullName))
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.Equal(t, metricName, resourceData.Get(ThresholdRuleFieldMetricName))
	require.Equal(t, window, resourceData.Get(ThresholdRuleFieldWindow))
	require.Equal(t, rollup, resourceData.Get(ThresholdRuleFieldRollup))
	require.Equal(t, string(aggregation), resourceData.Get(ThresholdRuleFieldAggregation))
	require.Equal(t, string(conditionOperator), resourceData.Get(ThresholdRuleFieldConditionOperator))
	require.Equal(t, conditionValue, resourceData.Get(ThresholdRuleFieldConditionValue))
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
	additionalAsserts(resourceData)
}

func TestShouldFailToUpdateTerraformStateForCustomEventSpecificationWithThresholdRuleWhenSeverityIsNotSupported(t *testing.T) {
	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:    restapi.ThresholdRuleType,
				Severity: 123,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithThresholdRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not a valid severity")
}

func TestShouldFailToUpdateTerraformStateForCustomEventSpecificationWithThresholdRuleWhenConditionOperatorTypeIsNotSupported(t *testing.T) {
	conditionOperator := "invalid"

	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          restapi.SeverityWarning.GetAPIRepresentation(),
				ConditionOperator: &conditionOperator,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithThresholdRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported condition operator")
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithThresholdRuleStateToDataModel(t *testing.T) {
	testMappingOfCustomEventSpecificationWithThresholdRuleTerraformStateToDataModel(t, func(resourceData *schema.ResourceData) { /* Default testcase without additional fields =< no additional mappings */
	}, func(spec restapi.CustomEventSpecification) { /* Default testcase without additional fields => no additional requires */
	})
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithThresholdRuleAndMetricPatternStateToDataModel(t *testing.T) {
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := restapi.MetricPatternOperatorTypeStartsWith

	additionalMappings := func(resourceData *schema.ResourceData) {
		resourceData.Set(ThresholdRuleFieldMetricPatternPrefix, prefix)
		resourceData.Set(ThresholdRuleFieldMetricPatternPostfix, postfix)
		resourceData.Set(ThresholdRuleFieldMetricPatternPlaceholder, placeholder)
		resourceData.Set(ThresholdRuleFieldMetricPatternOperator, operator)
	}

	additionalAsserts := func(spec restapi.CustomEventSpecification) {
		require.NotNil(t, spec.Rules[0].MetricPattern)
		require.Equal(t, prefix, spec.Rules[0].MetricPattern.Prefix)
		require.Equal(t, postfix, spec.Rules[0].MetricPattern.Postfix)
		require.Equal(t, placeholder, spec.Rules[0].MetricPattern.Placeholder)
		require.Equal(t, operator, spec.Rules[0].MetricPattern.Operator)
	}

	testMappingOfCustomEventSpecificationWithThresholdRuleTerraformStateToDataModel(t, additionalMappings, additionalAsserts)
}

func testMappingOfCustomEventSpecificationWithThresholdRuleTerraformStateToDataModel(t *testing.T, additionalMappings func(resourceData *schema.ResourceData), additionalAsserts func(spec restapi.CustomEventSpecification)) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	resourceData.Set(CustomEventSpecificationFieldFullName, resourceFullName)
	resourceData.Set(CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	resourceData.Set(CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	resourceData.Set(CustomEventSpecificationFieldTriggering, true)
	resourceData.Set(CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	resourceData.Set(CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	resourceData.Set(CustomEventSpecificationFieldEnabled, true)
	resourceData.Set(CustomEventSpecificationRuleSeverity, restapi.SeverityWarning.GetTerraformRepresentation())
	resourceData.Set(ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName)
	resourceData.Set(ThresholdRuleFieldWindow, customEventSpecificationWithThresholdRuleWindow)
	resourceData.Set(ThresholdRuleFieldRollup, customEventSpecificationWithThresholdRuleRollup)
	resourceData.Set(ThresholdRuleFieldAggregation, customEventSpecificationWithThresholdRuleAggregation)
	resourceData.Set(ThresholdRuleFieldConditionOperator, restapi.ConditionOperatorEquals.InstanaAPIValue())
	resourceData.Set(ThresholdRuleFieldConditionValue, customEventSpecificationWithThresholdRuleConditionValue)

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceFullName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, customEventSpecificationWithThresholdRuleMetricName, *customEventSpec.Rules[0].MetricName)
	require.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	require.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	require.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	require.Equal(t, restapi.ConditionOperatorEquals.InstanaAPIValue(), *customEventSpec.Rules[0].ConditionOperator)
	require.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithThresholdRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
}

func TestShouldFailToConvertCustomEventSpecificationWithThresholdRuleStateToDataModelWhenConditionOperationIsNotSupported(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, restapi.SeverityWarning.GetTerraformRepresentation())
	resourceData.Set(ThresholdRuleFieldConditionOperator, "invalid")

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a supported condition operator of the Instana Terraform provider")
}
