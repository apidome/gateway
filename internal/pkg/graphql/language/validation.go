package language

import (
	"fmt"

	"github.com/pkg/errors"
)

// Main validation function
func validateDocument(schema *document, docRoot *document) {

}

// spec.graphql.org/draft/#sec-Executable-Definitions
func validateExecutableDefinitions(doc document) {
	for _, def := range doc.Definitions() {
		_, isExecDef := def.(executableDefinition)

		if !isExecDef {
			panic("Document contains unexecutable definitions")
		}
	}
}

// http://spec.graphql.org/draft/#sec-Operation-Name-Uniqueness
func validateOperationNameUniqueness(doc document) {
	names := make(map[string]struct{})

	for _, op := range doc.Definitions() {
		opDef, isOpDef := op.(*operationDefinition)

		if isOpDef {
			_, nameExists := names[opDef.Name().Value()]

			if nameExists {
				panic("Operation name must be unique")
			}

			names[opDef.Name().Value()] = struct{}{}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Lone-Anonymous-Operation
func validateLoneAnonymousOperation(doc *document) {
	operations := getOperationDefinitions(doc)

	anonymous := getAnonymousOperationDefinitions(doc)

	if len(operations) > 1 {
		if len(anonymous) != 0 {
			panic("An anonymous operation must be the only operation in a document")
		}
	}
}

//! http://spec.graphql.org/draft/#sec-Single-root-field
func validateSingleRootField(doc *document, schema *schemaDefinition) {
	subscriptionOperations := getSubscriptionOperationDefinitions(doc)

	for _, sub := range subscriptionOperations {
		subscriptionType := getRootSubscriptionType(schema)
		selectionSet := sub.SelectionSet()

		variableValues := []variable{}
		groupedFieldSet := collectFields(doc, subscriptionType, selectionSet, variableValues)

		if len(groupedFieldSet) != 1 {
			panic(fmt.Sprintf("validateSingleRootField: groupFieldSet has more than 1 entry. subscriptionType: {{%v}}, selectionSet: {{%v}}, variableValues: {{%v}}",
				subscriptionType, selectionSet, variableValues))
		}
	}
}

//! http://spec.graphql.org/draft/#sec-Field-Selections-on-Objects-Interfaces-and-Unions-Types
func validateFieldSelectionsOnObjectsInterfaceAndUnionsTypes(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Field-Selection-Merging
func validateFieldSelectionMerging(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Leaf-Field-Selections
func validateLeafFieldSelections(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Argument-Names
func validateArgumentNames(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Argument-Uniqueness
func validateArgumentUniqueness(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Fragment-Name-Uniqueness
func validateFragmentNameUniqueness(doc document) {
	fragmentsPool := make(map[string]*fragmentDefinition)

	for _, def := range doc.definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			fragmentsPool[fragDef.fragmentName.value] = fragDef
		} else {
			panic(errors.New("Fragment names must be unique"))
		}
	}
}

//! http://spec.graphql.org/draft/#sec-Fragment-Spread-Type-Existence
func validateFragmentSpreadTypeExistence(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Fragments-On-Composite-Types
func validateFragmentsOnCompositeTypes(doc document) {

}

// Helper functions

func getOperationDefinitions(doc *document) []operationDefinition {
	operationDefinitions := make([]operationDefinition, 0)

	for _, def := range doc.Definitions() {
		opDef, isOpDef := def.(*operationDefinition)

		if isOpDef {
			operationDefinitions = append(operationDefinitions, *opDef)
		}
	}

	return operationDefinitions
}

func getExecutableDefinitions(doc *document) []executableDefinition {
	executableDefinitions := make([]executableDefinition, 0)

	for _, def := range doc.Definitions() {
		execDef, isExecDef := def.(executableDefinition)

		if isExecDef {
			executableDefinitions = append(executableDefinitions, execDef)
		}
	}

	return executableDefinitions
}

func getAnonymousOperationDefinitions(doc *document) []operationDefinition {
	anons := make([]operationDefinition, 0)

	for _, def := range doc.Definitions() {
		opDef, isOpDef := def.(*operationDefinition)

		if isOpDef {
			if opDef.Name() == nil {
				anons = append(anons, *opDef)
			}
		}
	}

	return anons
}

func getSubscriptionOperationDefinitions(doc *document) []operationDefinition {
	subscriptionOperationDefinitions := make([]operationDefinition, 0)

	for _, def := range doc.Definitions() {
		opDef, isOpDef := def.(*operationDefinition)

		if isOpDef {
			if opDef.OperationType() == operationSubscription {
				subscriptionOperationDefinitions = append(subscriptionOperationDefinitions, *opDef)
			}
		}
	}

	return subscriptionOperationDefinitions
}

func getRootSubscriptionType(schema *schemaDefinition) *rootOperationTypeDefinition {
	rootOperationTypeDefinitions := schema.RootOperationTypeDefinitions()
	for i := range rootOperationTypeDefinitions {
		if rootOperationTypeDefinitions[i].OperationType() == operationSubscription {
			return rootOperationTypeDefinitions[i]
		}
	}

	return nil
}
