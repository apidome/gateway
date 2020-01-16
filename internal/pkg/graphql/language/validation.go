package language

func validateDocument(sysDef typeSystemDefinition, docRoot *document) {

}

// spec.graphql.org/draft/#sec-Executable-Definitions
func validateExecutableDefinitions(doc document) {
	for _, def := range doc.Definitions {
		_, isExecDef := def.(executableDefinition)

		if !isExecDef {
			panic("Document contains unexecutable definitions")
		}
	}
}

// http://spec.graphql.org/draft/#sec-Operation-Name-Uniqueness
func validateOperationNameUniqueness(doc document) {
	names := make(map[string]struct{})

	for _, op := range doc.Definitions {
		opDef, isOpDef := op.(*operationDefinition)

		if isOpDef {
			_, nameExists := names[opDef.Name.Value]

			if nameExists {
				panic("Operation name must be unique")
			}

			names[opDef.Name.Value] = struct{}{}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Lone-Anonymous-Operation
// func validateLoneAnonymousOperation(doc document) {
// 	if len(doc.Definitions) > 1 {
// 		for _, op := range doc.Definitions {
// 			execDef, isExecDef := op.(executableDefinition)

// 			if isExecDef {

// 			}
// 		}
// 	}
// }

// http://spec.graphql.org/draft/#sec-Single-root-field
func validateSingleRootField(doc document) {

}

// http://spec.graphql.org/draft/#sec-Field-Selections-on-Objects-Interfaces-and-Unions-Types
func validateFieldSelectionsOnObjectsInterfaceAndUnionsTypes(doc document) {

}

// http://spec.graphql.org/draft/#sec-Field-Selection-Merging
func validateFieldSelectionMerging(doc document) {

}

// http://spec.graphql.org/draft/#sec-Leaf-Field-Selections
func validateLeafFieldSelections(doc document) {

}

// http://spec.graphql.org/draft/#sec-Argument-Names
func validateArgumentNames(doc document) {

}

// http://spec.graphql.org/draft/#sec-Argument-Uniqueness
func validateArgumentUniqueness(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragment-Name-Uniqueness
func validateFragmentNameUniqueness(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragment-Spread-Type-Existence
func validateFragmentSpreadTypeExistence(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragments-On-Composite-Types
func validateFragmentsOnCompositeTypes(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragments-Must-Be-Used
func validateFragmentsMustBeUsed(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragment-spread-target-defined
func validateFragmentSpreadTargetDefined(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragment-spreads-must-not-form-cycles
func validateFragmentSpreadsMustNotFormCycles(doc document) {

}

// http://spec.graphql.org/draft/#sec-Fragment-spread-is-possible
func validateFragmentSpreadIsPossible(doc document) {

}

// http://spec.graphql.org/draft/#sec-Values
func validateValuesOfCorrectType(doc document) {

}

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Names
func validateInputObjectFieldNames(doc document) {

}

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Uniqueness
func validateInputObjectFieldUniqueness(doc document) {

}

// http://spec.graphql.org/draft/#sec-Input-Object-Required-Fields
func validateInputObjectRequiredFieds(doc document) {

}

// http://spec.graphql.org/draft/#sec-Directives-Are-Defined
func validateDirectivesAreDefined(doc document) {

}

// http://spec.graphql.org/draft/#sec-Directives-Are-In-Valid-Locations
func validateDirectivesAreInValidLocations(doc document) {

}

// http://spec.graphql.org/draft/#sec-Directives-Are-Unique-Per-Location
func validateDirectivesAreUniquePerLocation(doc document) {

}

// http://spec.graphql.org/draft/#sec-Variable-Uniqueness
func validateVariableUniqueness(doc document) {

}

// http://spec.graphql.org/draft/#sec-Variables-Are-Input-Types
func validateVariableAreInputTypes(doc document) {

}

// http://spec.graphql.org/draft/#sec-All-Variable-Uses-Defined
func validateAllVariableUsesDefined(doc document) {

}

// http://spec.graphql.org/draft/#sec-All-Variables-Used
func validateAllVariablesUsed(doc document) {

}

// http://spec.graphql.org/draft/#sec-All-Variable-Usages-are-Allowed
func validateAllVariableUsagesAreAllowed(doc document) {

}
