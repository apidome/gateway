package language

import "github.com/pkg/errors"

// http://spec.graphql.org/draft/#sec-Fragments-Must-Be-Used
func validateFragmentsMustBeUsed(doc document) {
	fragmentSpreadTargets := make(map[string]struct{}, 0)

	// Extract fragment spread targets from all the definitions in the document.
	for _, def := range doc.Definitions {
		exeDef := def.(executableDefinition)
		extractUsedFragmentsNames(exeDef.GetSelectionSet(), fragmentSpreadTargets)
	}

	for _, def := range doc.Definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			if _, ok := fragmentSpreadTargets[fragDef.FragmentName.Value]; !ok {
				panic(errors.New("all defined fragments in a graphql document must be used"))
			}
		}
	}
}

func extractUsedFragmentsNames(selectionSet selectionSet, targetsSet map[string]struct{}) {
	for _, selection := range selectionSet {
		// If the selection is a fragment spread, append it's name to the names slice.
		// Else, extract the fragment names from all spreads in the selection's selectionSet.
		if fragmentSpread, ok := selection.(*fragmentSpread); ok {
			targetsSet[fragmentSpread.FragmentName.Value] = struct{}{}
		} else if selection.Selections() != nil {
			extractUsedFragmentsNames(*selection.Selections(), targetsSet)
		}
	}
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
func validateInputObjectRequiredFields(doc document) {

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
