package graphql_test

import (
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/testutil"
)

func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodIntValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: 2)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodBooleanValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: true)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodStringValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: "foo")
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodFloatValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: 1.1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_IntIntoFloat(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: 1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_IntIntoID(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: 1)
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_StringIntoID(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: "someIdString")
          }
        }
    `)
}
func TestValidate_ArgValuesOfCorrectType_ValidValue_GoodEnumValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: SIT)
          }
        }
    `)
}

func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_IntIntoString(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: 1)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringArg\" has invalid value 1.\n"+
					"Expected type \"String\", found 1.\n"+
					"Error: cannot parse *ast.IntValue to string",
				4, 39,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_FloatIntoString(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: 1.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringArg\" has invalid value 1.0.\n"+
					"Expected type \"String\", found 1.0.\n"+
					"Error: cannot parse *ast.FloatValue to string",
				4, 39,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_BooleanIntoString(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: true)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringArg\" has invalid value true.\n"+
					"Expected type \"String\", found true.\n"+
					"Error: cannot parse *ast.BooleanValue to string",
				4, 39,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidStringValues_UnquotedStringIntoString(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringArgField(stringArg: BAR)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringArg\" has invalid value BAR.\n"+
					"Expected type \"String\", found BAR.\n"+
					"Error: cannot parse *ast.EnumValue to string",
				4, 39,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_InvalidIntValues_StringIntoInt(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: "3")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"intArg\" has invalid value \"3\".\n"+
					"Expected type \"Int\", found \"3\".\n"+
					"Error: cannot parse *ast.StringValue to int",
				4, 33,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIntValues_BigIntIntoInt(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: 829384293849283498239482938)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"intArg\" has invalid value 829384293849283498239482938.\n"+
					"Expected type \"Int\", found 829384293849283498239482938.\n"+
					"Error: strconv.ParseInt: parsing \"829384293849283498239482938\": value out of range",
				4, 33,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIntValues_UnquotedStringIntoInt(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: FOO)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"intArg\" has invalid value FOO.\n"+
					"Expected type \"Int\", found FOO.\n"+
					"Error: cannot parse *ast.EnumValue to int",
				4, 33,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIntValues_SimpleFloatIntoInt(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: 3.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"intArg\" has invalid value 3.0.\n"+
					"Expected type \"Int\", found 3.0.\n"+
					"Error: cannot parse *ast.FloatValue to int",
				4, 33,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIntValues_FloatIntoInt(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            intArgField(intArg: 3.333)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"intArg\" has invalid value 3.333.\n"+
					"Expected type \"Int\", found 3.333.\n"+
					"Error: cannot parse *ast.FloatValue to int",
				4, 33,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_InvalidFloatValues_StringIntoFloat(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: "3.333")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"floatArg\" has invalid value \"3.333\".\n"+
					"Expected type \"Float\", found \"3.333\".\n"+
					"Error: cannot parse *ast.StringValue to float",
				4, 37,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidFloatValues_BooleanIntoFloat(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: true)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"floatArg\" has invalid value true.\n"+
					"Expected type \"Float\", found true.\n"+
					"Error: cannot parse *ast.BooleanValue to float",
				4, 37,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidFloatValues_UnquotedIntoFloat(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            floatArgField(floatArg: FOO)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"floatArg\" has invalid value FOO.\n"+
					"Expected type \"Float\", found FOO.\n"+
					"Error: cannot parse *ast.EnumValue to float",
				4, 37,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_InvalidBooleanValues_IntIntoBoolean(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: 2)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"booleanArg\" has invalid value 2.\n"+
					"Expected type \"Boolean\", found 2.\n"+
					"Error: cannot parse *ast.IntValue to bool",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidBooleanValues_FloatIntoBoolean(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: 1.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"booleanArg\" has invalid value 1.0.\n"+
					"Expected type \"Boolean\", found 1.0.\n"+
					"Error: cannot parse *ast.FloatValue to bool",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidBooleanValues_StringIntoBoolean(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: "true")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"booleanArg\" has invalid value \"true\".\n"+
					"Expected type \"Boolean\", found \"true\".\n"+
					"Error: cannot parse *ast.StringValue to bool",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidBooleanValues_UnquotedStringIntoBoolean(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            booleanArgField(booleanArg: TRUE)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"booleanArg\" has invalid value TRUE.\n"+
					"Expected type \"Boolean\", found TRUE.\n"+
					"Error: cannot parse *ast.EnumValue to bool",
				4, 41,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_InvalidIDValue_FloatIntoID(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: 1.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"idArg\" has invalid value 1.0.\n"+
					"Expected type \"ID\", found 1.0.\n"+
					"Error: cannot parse *ast.FloatValue to string",
				4, 31,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIDValue_BooleanIntoID(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: true)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"idArg\" has invalid value true.\n"+
					"Expected type \"ID\", found true.\n"+
					"Error: cannot parse *ast.BooleanValue to string",
				4, 31,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidIDValue_UnquotedIntoID(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            idArgField(idArg: SOMETHING)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"idArg\" has invalid value SOMETHING.\n"+
					"Expected type \"ID\", found SOMETHING.\n"+
					"Error: cannot parse *ast.EnumValue to string",
				4, 31,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_IntIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: 2)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value 2.\n"+
					"Expected type \"DogCommand\", found 2.\n"+
					"Error: Enum DogCommand cannot parse value: 2",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_FloatIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: 1.0)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value 1.0.\n"+
					"Expected type \"DogCommand\", found 1.0.\n"+
					"Error: Enum DogCommand cannot parse value: 1.0",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_StringIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: "SIT")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value \"SIT\".\n"+
					"Expected type \"DogCommand\", found \"SIT\".\n"+
					"Error: Enum DogCommand cannot parse value: SIT",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_BooleanIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: true)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value true.\n"+
					"Expected type \"DogCommand\", found true.\n"+
					"Error: Enum DogCommand cannot parse value: true",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_UnknownEnumValueIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: JUGGLE)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value JUGGLE.\n"+
					"Expected type \"DogCommand\", found JUGGLE.\n"+
					"Error: Enum DogCommand cannot parse value: JUGGLE",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidEnumValue_DifferentCaseEnumValueIntoEnum(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            doesKnowCommand(dogCommand: sit)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"dogCommand\" has invalid value sit.\n"+
					"Expected type \"DogCommand\", found sit.\n"+
					"Error: Enum DogCommand cannot parse value: sit",
				4, 41,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_ValidListValue_GoodListValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringListArgField(stringListArg: ["one", "two"])
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidListValue_EmptyListValue(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringListArgField(stringListArg: [])
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidListValue_SingleValueIntoList(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringListArgField(stringListArg: "one")
          }
        }
        `)
}

func TestValidate_ArgValuesOfCorrectType_InvalidListValue_IncorrectItemType(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringListArgField(stringListArg: ["one", 2])
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringListArg\" has invalid value [\"one\", 2].\n"+
					"In element #2: Expected type \"String\", found 2.\n"+
					"In element #2: Error: cannot parse *ast.IntValue to string",
				4, 47,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidListValue_SingleValueOfIncorrentType(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            stringListArgField(stringListArg: 1)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"stringListArg\" has invalid value 1.\n"+
					"Expected type \"String\", found 1.\n"+
					"Error: cannot parse *ast.IntValue to string",
				4, 47,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_ArgOnOptionalArg(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            isHousetrained(atOtherHomes: true)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_NoArgOnOptionalArg(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog {
            isHousetrained
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_MultipleArgs(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleReqs(req1: 1, req2: 2)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_MultipleArgsReverseOrder(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleReqs(req2: 2, req1: 1)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_NoArgsOnMultipleOptional(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOpts
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_OneArgOnMultipleOptional(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOpts(opt1: 1)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_SecondArgOnMultipleOptional(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOpts(opt2: 1)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_MultipleRequiredsOnMixedList(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOptAndReq(req1: 3, req2: 4)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_MultipleRequiredsAndOptionalOnMixedList(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOptAndReq(req1: 3, req2: 4, opt1: 5)
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidNonNullableValue_AllRequiredsAndOptionalOnMixedList(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleOptAndReq(req1: 3, req2: 4, opt1: 5, opt2: 6)
          }
        }
        `)
}

func TestValidate_ArgValuesOfCorrectType_InvalidNonNullableValue_IncorrectValueType(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleReqs(req2: "two", req1: "one")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"req2\" has invalid value \"two\".\n"+
					"Expected type \"Int\", found \"two\".\n"+
					"Error: cannot parse *ast.StringValue to int",
				4, 32,
			),
			testutil.RuleError(
				"Argument \"req1\" has invalid value \"one\".\n"+
					"Expected type \"Int\", found \"one\".\n"+
					"Error: cannot parse *ast.StringValue to int",
				4, 45,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidNonNullableValue_IncorrectValueAndMissingArgument(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            multipleReqs(req1: "one")
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"req1\" has invalid value \"one\".\n"+
					"Expected type \"Int\", found \"one\".\n"+
					"Error: cannot parse *ast.StringValue to int",
				4, 32,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_OptionalArg_DespiteRequiredFieldInType(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_PartialObject_OnlyRequired(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: { requiredField: true })
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_PartialObject_RequiredFieldCanBeFalsey(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: { requiredField: false })
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_PartialObject_IncludingRequired(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
			  complexArgField(complexArg: { requiredField: false, intField: 4 })
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_FullObject(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: {
              requiredField: true,
              intField: 4,
              stringField: "foo",
              booleanField: false,
              stringListField: ["one", "two"]
            })
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_ValidInputObjectValue_FullObject_WithFieldsInDifferentOrder(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: {
              stringListField: ["one", "two"],
              booleanField: false,
              requiredField: true,
              stringField: "foo",
              intField: 4,
            })
          }
        }
        `)
}

func TestValidate_ArgValuesOfCorrectType_InvalidInputObjectValue_PartialObject_MissingRequired(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: { intField: 4 })
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"complexArg\" has invalid value {intField: 4}.\nIn field \"requiredField\": Expected \"Boolean!\", found null.",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidInputObjectValue_PartialObject_InvalidFieldType(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: {
              stringListField: ["one", 2],
              requiredField: true,
            })
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"complexArg\" has invalid value {stringListField: [\"one\", 2], requiredField: true}.\n"+
					"In field \"stringListField\": In element #2: Expected type \"String\", found 2.\n"+
					"In field \"stringListField\": In element #2: Error: cannot parse *ast.IntValue to string",
				4, 41,
			),
		})
}
func TestValidate_ArgValuesOfCorrectType_InvalidInputObjectValue_PartialObject_UnknownFieldArg(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          complicatedArgs {
            complexArgField(complexArg: {
              requiredField: true,
              unknownField: "value"
            })
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				"Argument \"complexArg\" has invalid value {requiredField: true, unknownField: \"value\"}.\nIn field \"unknownField\": Unknown field.",
				4, 41,
			),
		})
}

func TestValidate_ArgValuesOfCorrectType_DirectiveArguments_WithDirectivesOfValidType(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog @include(if: true) {
            name
          }
          human @skip(if: false) {
            name
          }
        }
        `)
}
func TestValidate_ArgValuesOfCorrectType_DirectiveArguments_WithDirectivesWithIncorrectTypes(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.ArgumentsOfCorrectTypeRule, `
        {
          dog @include(if: "yes") {
            name @skip(if: ENUM)
          }
        }
        `,
		[]gqlerrors.FormattedError{
			testutil.RuleError(
				`Argument "if" has invalid value "yes".`+
					"\nExpected type \"Boolean\", found \"yes\"."+
					"\nError: cannot parse *ast.StringValue to bool",
				3, 28,
			),
			testutil.RuleError(
				`Argument "if" has invalid value ENUM.`+
					"\nExpected type \"Boolean\", found ENUM."+
					"\nError: cannot parse *ast.EnumValue to bool",
				4, 28,
			),
		})
}
