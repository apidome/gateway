package language

func validateDocument(sysDef typeSystemDefinition, docRoot document) {

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

}

// http://spec.graphql.org/draft/#sec-Lone-Anonymous-Operation
func validateLoneAnonymousOperaation(doc document) {

}

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
