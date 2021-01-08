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

//! http://spec.graphql.org/draft/#sec-Field-Selection-Merging
func validateFieldSelectionMerging(doc document) {

}

//! http://spec.graphql.org/draft/#sec-Leaf-Field-Selections
func validateLeafFieldSelections(schema, doc document) {
	errMsg := "Field selections on scalars or enums are never allowed\n" +
	" because they are the leaf nodes of any GraphQL querys"

	for _, def := range doc.definitions {
		if execDef, isExecDef := def.(executableDefinition); isExecDef {
			if execDef.SelectionSet() != nil {
				if !isLeafSelectionValid(
					schema,
					execDef.SelectionSet(),
					execDef.SelectionSet(),
					getRootQueryTypeDefinition(&schema),
					getFragmentsPool(&doc),
				 ) {
					panic(errors.New(errMsg))
				}
			}
		}
	}
}

func isLeafSelectionValid(
	schema document,
	rootSelectionSet selectionSet,
	selectionSet selectionSet,
	parentType typeDefinition,
	fragmentsPool map[string]*fragmentDefinition,
) bool {
	for _, selection := range selectionSet {
		nextSelectionSet := *selection.SelectionSet()
		if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			nextSelectionSet = fragmentsPool[fragSpread.fragmentName.value].selectionSet
		}

		// If the selection have no sub selection, it a leaf selection.
		isLeafSelection := nextSelectionSet == nil

		// Get the field definition from the schems.
		fieldDef := getFieldDefinitionByFieldSelection(
			parentType,
			selection,
			selectionSet,
			&schema,
			fragmentsPool,
		)

		// Get the type definition of the selection's return value.
		selectionType := getTypeDefinitionByType(&schema, fieldDef.Type())

		// If selectionType is a scalar or enum:
		// 	The subselection set of that selection must be empty
		// If selectionType is an interface, union, or object
		// 	The subselection set of that selection must NOT BE empty
		switch selectionType.(type) {
		case *scalarTypeDefinition, *enumTypeDefinition:
			if !isLeafSelection {
				return false
			}
		case *interfaceTypeDefinition, *unionTypeDefinition, *objectTypeDefinition:
			if isLeafSelection {
				return false
			}
		}

		if !isLeafSelection {
			if !isLeafSelectionValid(
				schema,
				rootSelectionSet,
				nextSelectionSet,
				selectionType,
				fragmentsPool,
			) {
				return false
			}
		}
	}

	return true
}

//! http://spec.graphql.org/draft/#sec-Argument-Names
func validateArgumentNames(schema, doc document) {
	fragmentsPool := getFragmentsPool(&doc)

	for _, def := range doc.Definitions() {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			checkArgumentNamesInOpDef(
				schema,
				opDef.selectionSet,
				opDef.selectionSet,
				fragmentsPool,
			)
		}
	}
}

func checkArgumentNamesInOpDef(
	schema document,
	rootSelection selectionSet,
	selectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) {
	for _, selection := range selectionSet {
		if selection.Directives() != nil {
			for _, dir := range *selection.Directives() {
				if dir.arguments != nil {
					for _, arg := range *dir.arguments {
						getDirectiveArgumentDefinition(&schema, dir, arg)
					}
				}
			}
		}

		switch s := selection.(type) {
		case *field:
			if s.arguments != nil {
				for _, arg := range *s.arguments {
					getFieldArgumentDefinition(&schema, rootSelection, s, arg, fragmentsPool)
				}
			}
		case *inlineFragment:
			checkArgumentNamesInOpDef(
				schema,
				rootSelection,
				s.selectionSet,
				fragmentsPool,
			)
		case *fragmentSpread:
			checkArgumentNamesInOpDef(
				schema,
				rootSelection,
				fragmentsPool[s.fragmentName.value].selectionSet,
				fragmentsPool,
			)
		}
	}
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
func validateFragmentSpreadTypeExistence(schema, doc document) {
	for _, def := range doc.Definitions() {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			doesFragmentSpreadTypeExist(
				schema,
				exeDef.SelectionSet(),
				getFragmentsPool(&doc),
			)			
		}
	}
}

func doesFragmentSpreadTypeExist(
	schema document,
	selectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) {
	for _, selection := range selectionSet {
		nextSelectionSet := *(selection.SelectionSet())
		switch t := selection.(type) {
		case *fragmentSpread:
			fragment := fragmentsPool[t.fragmentName.value]
			fragmentType := fragment.typeCondition.namedType._type()
			getTypeDefinitionByType(&schema, fragmentType)
			nextSelectionSet = fragment.SelectionSet()
		case *inlineFragment:
			inlineFragType := t.typeCondition.namedType._type()
			getTypeDefinitionByType(&schema, inlineFragType)
		}

		doesFragmentSpreadTypeExist(
			schema,
			nextSelectionSet,
			fragmentsPool,
		)
	}
}

//! http://spec.graphql.org/draft/#sec-Fragments-On-Composite-Types
func validateFragmentsOnCompositeTypes(schema, doc document) {
	errMsg := "Fragments can only be declared on unions, interfaces, and objects." +
	" They are invalid on scalars. They can only be applied on non‐leaf fields." +
	" This rule applies to both inline and named fragments."

	for _, def := range doc.Definitions() {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			if !isFragmentOnCompositeType(
				schema,
				exeDef.SelectionSet(),
				getFragmentsPool(&doc),
			) {
				panic(errors.New(errMsg))
			}
		}
	}
}

func isFragmentOnCompositeType(
	schema document,
	selectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) bool {
	for _, selection := range selectionSet {
		nextSelectionSet := *(selection.SelectionSet())
		switch t := selection.(type) {
		case *fragmentSpread:
			fragment := fragmentsPool[t.fragmentName.value]
			nextSelectionSet = fragment.SelectionSet()
			fragmentType := fragment.typeCondition.namedType._type()
			typeDef := getTypeDefinitionByType(&schema, fragmentType)
			if !isCompositeType(typeDef) {
				return false
			}
		case *inlineFragment:
			inlineFragType := t.typeCondition.namedType._type()
			typeDef := getTypeDefinitionByType(&schema, inlineFragType)
			if !isCompositeType(typeDef) {
				return false
			}
		}

		if !isFragmentOnCompositeType(schema, nextSelectionSet, fragmentsPool) {
			return false
		}
	}

	return true
}

func isCompositeType(typeDef typeDefinition) bool {
	switch typeDef.(type) {
	case *objectTypeDefinition, *interfaceTypeDefinition, *unionTypeDefinition:
		return true
	}

	return false
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
