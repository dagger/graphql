package graphql_test

import (
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/testutil"
)

func TestValidate_FieldsOnCorrectType_ObjectFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment objectFieldSelection on Dog {
        __typename
        name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_AliasedObjectFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedObjectFieldSelection on Dog {
        tn : __typename
        otherName : name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_InterfaceFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment interfaceFieldSelection on Pet {
        __typename
        name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_AliasedInterfaceFieldSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment interfaceFieldSelection on Pet {
        otherName : name
      }
    `)
}
func TestValidate_FieldsOnCorrectType_LyingAliasSelection(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment lyingAliasSelection on Dog {
        name : nickname
      }
    `)
}
func TestValidate_FieldsOnCorrectType_IgnoresFieldsOnUnknownType(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment unknownSelection on UnknownType {
        unknownField
      }
    `)
}
func TestValidate_FieldsOnCorrectType_ReportErrorsWhenTheTypeIsKnownAgain(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment typeKnownAgain on Pet {
        unknown_pet_field {
          ... on Cat {
            unknown_cat_field
          }
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "unknown_pet_field" on type "Pet".`, 3, 9),
		testutil.RuleError(`Cannot query field "unknown_cat_field" on type "Cat".`, 5, 13),
	})
}
func TestValidate_FieldsOnCorrectType_FieldNotDefinedOnFragment(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment fieldNotDefined on Dog {
        meowVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "meowVolume" on type "Dog". Did you mean "barkVolume"?`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_IgnoreDeeplyUnknownField(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment deepFieldNotDefined on Dog {
        unknown_field {
          deeper_unknown_field
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "unknown_field" on type "Dog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_SubFieldNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment subFieldNotDefined on Human {
        pets {
          unknown_field
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "unknown_field" on type "Pet".`, 4, 11),
	})
}
func TestValidate_FieldsOnCorrectType_FieldNotDefinedOnInlineFragment(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment fieldNotDefined on Pet {
        ... on Dog {
          meowVolume
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "meowVolume" on type "Dog". Did you mean "barkVolume"?`, 4, 11),
	})
}
func TestValidate_FieldsOnCorrectType_AliasedFieldTargetNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedFieldTargetNotDefined on Dog {
        volume : mooVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "mooVolume" on type "Dog". Did you mean "barkVolume"?`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_AliasedLyingFieldTargetNotDefined(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment aliasedLyingFieldTargetNotDefined on Dog {
        barkVolume : kawVolume
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "kawVolume" on type "Dog". Did you mean "barkVolume"?`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_NotDefinedOnInterface(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment notDefinedOnInterface on Pet {
        tailLength
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "tailLength" on type "Pet".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_DefinedOnImplementorsButNotOnInterface(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment definedOnImplementorsButNotInterface on Pet {
        nickname
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "nickname" on type "Pet". Did you mean to use an inline fragment on "Cat" or "Dog"?`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_MetaFieldSelectionOnUnion(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment directFieldSelectionOnUnion on CatOrDog {
        __typename
      }
    `)
}
func TestValidate_FieldsOnCorrectType_DirectFieldSelectionOnUnion(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment directFieldSelectionOnUnion on CatOrDog {
        directField
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "directField" on type "CatOrDog".`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_DefinedImplementorsQueriedOnUnion(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment definedOnImplementorsQueriedOnUnion on CatOrDog {
        name
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Cannot query field "name" on type "CatOrDog". Did you mean to use an inline fragment on "Being", "Pet", "Canine", "Cat", or "Dog"?`, 3, 9),
	})
}
func TestValidate_FieldsOnCorrectType_ValidFieldInInlineFragment(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `
      fragment objectFieldSelection on Pet {
        ... on Dog {
          name
        }
        ... {
          name
        }
      }
    `)
}

func TestValidate_FieldsOnCorrectTypeErrorMessage_WorksWithNoSuggestions(t *testing.T) {
	message := graphql.UndefinedFieldMessage("f", "T", []string{}, []string{})
	expected := `Cannot query field "f" on type "T".`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}

func TestValidate_FieldsOnCorrectTypeErrorMessage_WorksWithNoSmallNumbersOfTypeSuggestions(t *testing.T) {
	message := graphql.UndefinedFieldMessage("f", "T", []string{"A", "B"}, []string{})
	expected := `Cannot query field "f" on type "T". ` +
		`Did you mean to use an inline fragment on "A" or "B"?`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}

func TestValidate_FieldsOnCorrectTypeErrorMessage_WorksWithNoSmallNumbersOfFieldSuggestions(t *testing.T) {
	message := graphql.UndefinedFieldMessage("f", "T", []string{}, []string{"z", "y"})
	expected := `Cannot query field "f" on type "T". ` +
		`Did you mean "z" or "y"?`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}
func TestValidate_FieldsOnCorrectTypeErrorMessage_OnlyShowsOneSetOfSuggestionsAtATimePreferringTypes(t *testing.T) {
	message := graphql.UndefinedFieldMessage("f", "T", []string{"A", "B"}, []string{"z", "y"})
	expected := `Cannot query field "f" on type "T". ` +
		`Did you mean to use an inline fragment on "A" or "B"?`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}

func TestValidate_FieldsOnCorrectTypeErrorMessage_LimitLotsOfTypeSuggestions(t *testing.T) {
	message := graphql.UndefinedFieldMessage("f", "T", []string{"A", "B", "C", "D", "E", "F"}, []string{})
	expected := `Cannot query field "f" on type "T". ` +
		`Did you mean to use an inline fragment on "A", "B", "C", "D", or "E"?`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}

func TestValidate_FieldsOnCorrectTypeErrorMessage_LimitLotsOfFieldSuggestions(t *testing.T) {
	message := graphql.UndefinedFieldMessage(
		"f", "T", []string{}, []string{"z", "y", "x", "w", "v", "u"},
	)
	expected := `Cannot query field "f" on type "T". ` +
		`Did you mean "z", "y", "x", "w", or "v"?`
	if message != expected {
		t.Fatalf("Unexpected message, expected: %v, got %v", expected, message)
	}
}

func TestValidate_FieldsOnCorrectType_NilCrash(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.FieldsOnCorrectTypeRule, `mutation{o}`)
}
