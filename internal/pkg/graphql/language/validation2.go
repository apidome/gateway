package language

import (
	"github.com/pkg/errors"
)

// http://spec.graphql.org/draft/#sec-Fragments-Must-Be-Used
func validateFragmentsMustBeUsed(doc document) {
	fragmentSpreadTargets := make(map[string]struct{}, 0)

	// Extract fragment spread targets from all the definitions in the document.
	for _, def := range doc.Definitions {
		exeDef := def.(executableDefinition)
		extractUsedFragmentsNames(exeDef.GetSelectionSet(), fragmentSpreadTargets)
	}

	// For each fragment definition in the document, check if it's name exists
	// in the fragment spread targets set. If not, panic.
	for _, def := range doc.Definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			if _, ok := fragmentSpreadTargets[fragDef.FragmentName.Value]; !ok {
				panic(errors.New("Defined fragments must be used within " +
					"a document"))
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Fragment-spread-target-defined
func validateFragmentSpreadTargetDefined(doc document) {
	fragmentSpreadTargets := make(map[string]struct{}, 0)
	fragmentDefinitionsNames := make(map[string]struct{})

	// Collect all the fragment definitions' names into a dictionary of names.
	for _, def := range doc.Definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			fragmentDefinitionsNames[fragDef.FragmentName.Value] = struct{}{}
		}
	}

	// Extract fragment spread targets from all the definitions in the document.
	for _, def := range doc.Definitions {
		exeDef := def.(executableDefinition)
		extractUsedFragmentsNames(exeDef.GetSelectionSet(), fragmentSpreadTargets)
	}

	// Verify that all the targets in the targets set are referring to a defined fragment.
	// If not, panic.
	for target := range fragmentSpreadTargets {
		if _, ok := fragmentDefinitionsNames[target]; !ok {
			panic(errors.New("Named fragment spreads must refer to fragments" +
				" defined within the document"))
		}
	}
}

// http://spec.graphql.org/draft/#sec-Fragment-spreads-must-not-form-cycles
func validateFragmentSpreadsMustNotFormCycles(doc document) {
	fragmentsPool := getFragmentsPool(doc)

	// For each fragment, call detectFragmentCycles()
	for _, fragDef := range fragmentsPool {
		visited := make(map[string]struct{}, 0)
		detectFragmentCycles(*fragDef, visited, fragmentsPool)
	}
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
func validateDirectivesAreDefined(schema, doc document) {
	for _, def := range doc.Definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			if !checkIfDirectivesInOperationDefined(schema, exeDef.GetSelectionSet()) {
				panic(errors.New("GraphQL servers define what directives they " +
					"support. For each usage of a directive, the directive must be " +
					"available on that server"))
			}
		}
	}
}

func checkIfDirectivesInOperationDefined(schema document, selectionSet selectionSet) bool {
	// A map that contains all graphql's built in directives for quick access.
	builtInDirectiveNames := map[string]struct{}{
		"skip":       struct{}{},
		"include":    struct{}{},
		"deprecated": struct{}{},
	}

	// For each selection, look for directives.
	for _, selection := range selectionSet {
		// If the selection is a field
		if field, isField := selection.(*field); isField {
			// And the field contains directives
			if field.Directives != nil {
				// For each directive, decide whether it is defined or not.
				for _, directive := range *field.Directives {
					// A flag that indicated if the directive is defined or not.
					var isDirectiveDefined bool

					// Check if the directive is a build in directive. If it is, turn on the
					// flag. Else, Search the directive in the schema.
					if _, isBuiltInDirective := builtInDirectiveNames[directive.Name.Value]; isBuiltInDirective {
						isDirectiveDefined = true
					} else {
						// For each definition in the schema, if it is a directive definition,
						// compare its name to the current directive from the document.
						for _, def := range schema.Definitions {
							if directiveDef, isDirectiveDef := def.(*directiveDefinition); isDirectiveDef {
								// If the name are equal, turn on the flag.
								if directiveDef.Name.Value == directive.Name.Value {
									isDirectiveDefined = true
								}
							}
						}
					}

					// If the flag is off (which means the directive is not a built in
					// directive and we could not find a proper directive definition in
					// the schema), return false.
					if !isDirectiveDefined {
						return false
					}
				}
			}

			// If the field contains a selection set, check it's directives too.
			if field.SelectionSet != nil {
				checkIfDirectivesInOperationDefined(schema, *field.SelectionSet)
			}
			// If the selection is an inline fragment, check it's selection set's directives.
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			checkIfDirectivesInOperationDefined(schema, inlineFrag.SelectionSet)
		}
	}

	// If we arrived this line, it means we did not return false for any directive,
	// so we return true.
	return true
}

// http://spec.graphql.org/draft/#sec-Directives-Are-In-Valid-Locations
func validateDirectivesAreInValidLocations(doc document) {

}

// http://spec.graphql.org/draft/#sec-Directives-Are-Unique-Per-Location
func validateDirectivesAreUniquePerLocation(doc document) {

}

// http://spec.graphql.org/draft/#sec-Variable-Uniqueness
func validateVariableUniqueness(doc document) {
	for _, def := range doc.Definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.VariableDefinitions != nil {
				variablesSet := make(map[string]struct{}, 0)

				for _, variable := range *opDef.VariableDefinitions {
					if _, isVariableExists := variablesSet[variable.Variable.Name.Value]; isVariableExists {
						panic(errors.New("If any operation defines more than one" +
							" variable with the same name, it is ambiguous and invalid"))
					}

					variablesSet[variable.Variable.Name.Value] = struct{}{}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Variables-Are-Input-Types
func validateVariableAreInputTypes(schema document, doc document) {
	for _, def := range doc.Definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.VariableDefinitions != nil {
				for _, variable := range *opDef.VariableDefinitions {
					if !isInputType(schema, variable.Type) {
						panic(errors.New("Variables can only be input types. " +
							"Objects, unions, and interfaces cannot be used as inputs."))
					}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-All-Variable-Uses-Defined
func validateAllVariableUsesDefined(doc document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.Definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			usedVariablesSet := make(map[string]struct{}, 0)
			extractUsedVariablesNames(opDef.SelectionSet, usedVariablesSet, fragmentsPool)

			definedVariables := make(map[string]struct{}, 0)
			if opDef.VariableDefinitions != nil {
				for _, variable := range *opDef.VariableDefinitions {
					definedVariables[variable.Variable.Name.Value] = struct{}{}
				}
			}

			for usedVariable := range usedVariablesSet {
				if _, isUsedVariableDefined := definedVariables[usedVariable]; !isUsedVariableDefined {
					panic(errors.New("Variables are scoped on a per‐operation " +
						"basis. That means that any variable used within the context of " +
						"an operation must be defined at the top level of that operation"))
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-All-Variables-Used
func validateAllVariablesUsed(doc document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.Definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.VariableDefinitions != nil {
				usedVariablesSet := make(map[string]struct{}, 0)
				extractUsedVariablesNames(opDef.SelectionSet, usedVariablesSet, fragmentsPool)

				for _, variable := range *opDef.VariableDefinitions {
					if _, isVariableUsed := usedVariablesSet[variable.Variable.Name.Value]; !isVariableUsed {
						panic(errors.New("All variables defined by an operation" +
							" must be used in that operation or a fragment transitively " +
							"included by that operation"))
					}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-All-Variable-Usages-are-Allowed
func validateAllVariableUsagesAreAllowed(schema, doc document) {
	for _, def := range doc.Definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			variableUsages := make(map[interface{}]*variable)
			fragmentsPool := getFragmentsPool(doc)
			collectVariableUsages(opDef.SelectionSet, variableUsages, fragmentsPool)

			for varUsage := range variableUsages {
				for _, varDef := range *opDef.VariableDefinitions {
					if varUsage == varDef.Variable.Name.Value {
						if !isVariableUsageAllowed(&varDef) {
							panic(errors.New("Variable usages must be compatible" +
								" with the arguments they are passed to"))
						}
					}
				}
			}
		}
	}
}

func isVariableUsageAllowed(varDef *variableDefinition /*, variableUsage */) bool {
	var (
		hasNonNullVariableDefaultValue bool
		hasLocationDefaultValue        bool
	)

	// Let locationType be the expected type of the Argument, ObjectField, or ListValue
	// entry where variableUsage is located.
	// TODO: Replace this with a function argument.
	locationType := (&nonNullType{})._type()

	// Let variableType be the expected type of variableDefinition
	variableType := varDef.Type

	_, isVariableTypeANonNullType := variableType.(*nonNullType)
	nonNullLocationType, isLocationTypeANonNullType := locationType.(*nonNullType)

	// If locationType is a non‐null type AND variableType is NOT a non‐null type:
	if isLocationTypeANonNullType && !isVariableTypeANonNullType {
		// Let hasNonNullVariableDefaultValue be true if a default value exists for
		// variableDefinition and is not the value null.
		if varDef.DefaultValue != nil {
			if _, isNullValue := (*varDef).DefaultValue.Value.(*nullValue); !isNullValue {
				hasNonNullVariableDefaultValue = true
			}
		}

		// Let hasLocationDefaultValue be true if a default value exists for the Argument
		// or ObjectField where variableUsage is located.
		//if locationType.DefaultValue {
		//	hasLocationDefaultValue = true
		//}

		// If hasNonNullVariableDefaultValue is NOT true AND hasLocationDefaultValue is
		// NOT true, return false.
		if !hasNonNullVariableDefaultValue && !hasLocationDefaultValue {
			return false
		}

		// Let nullableLocationType be the unwrapped nullable type of locationType.
		nullableLocationType := nonNullLocationType.OfType

		// Check if the types are compatible.
		return areTypesCompatible(variableType, nullableLocationType)
	}

	// Check if the types are compatible.
	return areTypesCompatible(variableType, locationType)
}

func areTypesCompatible(variableType, locationType _type) bool {
	nonNullVariableType, isVariableTypeIsNonNullType := variableType.(*nonNullType)
	nonNullLocationType, isLocationTypeIsNonNullType := locationType.(*nonNullType)

	// If locationType is a non‐null type:
	if isLocationTypeIsNonNullType {
		// If variableType is NOT a non‐null type, return false.
		if !isVariableTypeIsNonNullType {
			return false
		}

		// Let nullableLocationType be the unwrapped nullable type of locationType.
		nullableLocationType := nonNullLocationType.OfType

		// Let nullableVariableType be the unwrapped nullable type of variableType.
		nullableVariableType := nonNullVariableType.OfType

		// Return the result of areTypesCompatible with the unwrapped types.
		return areTypesCompatible(nullableVariableType, nullableLocationType)
	}

	// Otherwise, if variableType is a non‐null type:
	if isVariableTypeIsNonNullType {
		// Let nullableVariableType be the unwrapped nullable type of variableType.
		nullableVariableType := nonNullVariableType.OfType

		// Return the result of areTypesCompatible with the unwrapped types.
		return areTypesCompatible(nullableVariableType, locationType)
	}

	listVariableType, isVariableTypeAListType := variableType.(*listType)
	listLocationType, isLocationTypeAListType := locationType.(*listType)

	// Otherwise, if locationType is a list type:
	if isLocationTypeAListType {
		// If variableType is NOT a list type, return false.
		if !isVariableTypeAListType {
			return false
		}

		// Let itemLocationType be the unwrapped item type of locationType.
		itemLocationType := listLocationType.OfType

		// Let itemVariableType be the unwrapped item type of variableType.
		itemVariableType := listVariableType.OfType

		// Return the result of areTypesCompatible with the unwrapped types.
		return areTypesCompatible(itemVariableType, itemLocationType)
	}

	// Return true if variableType and locationType are identical, otherwise false.
	return variableType.GetTypeName() == locationType.GetTypeName()
}

func collectVariableUsages(selectionSet selectionSet,
	usages map[interface{}]*variable,
	fragmentsPool map[string]*fragmentDefinition) {
	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection is a field, extract its variable usages
		// (the field's arguments, and the field's directive's arguments).
		if field, isField := selection.(*field); isField {
			if field.Arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.Arguments {
					if _var, isVariable := arg.Value.(*variable); isVariable {
						usages[&arg] = _var
					} else if list, isList := arg.Value.(*listValue); isList {
						for _, item := range list.Values {
							if _var, isVariable := item.(*variable); isVariable {
								usages[&list] = _var
							}
						}
					} else if object, isObject := arg.Value.(*objectValue); isObject {
						for _, field := range object.Values {
							if _var, isVariable := field.Value.(*variable); isVariable {
								usages[&field] = _var
							}
						}
					}
				}
			}

			// If the field has directives, check their arguments too.
			if field.Directives != nil {
				for _, directive := range *field.Directives {
					if directive.Arguments != nil {
						// For each argument, if it is a variable, add it to the variables set.
						// If it is a list, check each item in the list.
						// If it is an object, check each value in the object.
						for _, arg := range *directive.Arguments {
							if _var, isVariable := arg.Value.(*variable); isVariable {
								usages[&arg] = _var
							} else if list, isList := arg.Value.(*listValue); isList {
								for _, item := range list.Values {
									if _var, isVariable := item.(*variable); isVariable {
										usages[&list] = _var
									}
								}
							} else if object, isObject := arg.Value.(*objectValue); isObject {
								for _, field := range object.Values {
									if _var, isVariable := field.Value.(*variable); isVariable {
										usages[&field] = _var
									}
								}
							}
						}
					}
				}
			}

			// If the field has a selection set, extract its variables too.
			if field.SelectionSet != nil {
				collectVariableUsages(*field.SelectionSet, usages, fragmentsPool)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			collectVariableUsages(*inlineFrag.GetSelections(), usages, fragmentsPool)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			collectVariableUsages(fragmentsPool[fragSpread.FragmentName.Value].SelectionSet,
				usages,
				fragmentsPool)
		}
	}
}

func extractUsedFragmentsNames(selectionSet selectionSet, targetsSet map[string]struct{}) {
	for _, selection := range selectionSet {
		// If the selection is a fragment spread, append it's name to the names slice.
		// Else, extract the fragment names from all spreads in the selection's selectionSet.
		if fragmentSpread, ok := selection.(*fragmentSpread); ok {
			targetsSet[fragmentSpread.FragmentName.Value] = struct{}{}
		} else if selection.GetSelections() != nil {
			extractUsedFragmentsNames(*selection.GetSelections(), targetsSet)
		}
	}
}

func extractUsedVariablesNames(selectionSet selectionSet,
	variablesSet map[string]struct{},
	fragmentsPool map[string]*fragmentDefinition) {
	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection is a field, extract its variable usages
		// (the field's arguments, and the field's directive's arguments).
		if field, isField := selection.(*field); isField {
			if field.Arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.Arguments {
					if _var, isVariable := arg.Value.(*variable); isVariable {
						variablesSet[_var.Name.Value] = struct{}{}
					} else if list, isList := arg.Value.(*listValue); isList {
						for _, item := range list.Values {
							if _var, isVariable := item.(*variable); isVariable {
								variablesSet[_var.Name.Value] = struct{}{}
							}
						}
					} else if object, isObject := arg.Value.(*objectValue); isObject {
						for _, field := range object.Values {
							if _var, isVariable := field.Value.(*variable); isVariable {
								variablesSet[_var.Name.Value] = struct{}{}
							}
						}
					}
				}
			}

			// If the field has directives, check their arguments too.
			if field.Directives != nil {
				for _, directive := range *field.Directives {
					if directive.Arguments != nil {
						// For each argument, if it is a variable, add it to the variables set.
						// If it is a list, check each item in the list.
						// If it is an object, check each value in the object.
						for _, arg := range *directive.Arguments {
							if _var, isVariable := arg.Value.(*variable); isVariable {
								variablesSet[_var.Name.Value] = struct{}{}
							} else if list, isList := arg.Value.(*listValue); isList {
								for _, item := range list.Values {
									if _var, isVariable := item.(*variable); isVariable {
										variablesSet[_var.Name.Value] = struct{}{}
									}
								}
							} else if object, isObject := arg.Value.(*objectValue); isObject {
								for _, field := range object.Values {
									if _var, isVariable := field.Value.(*variable); isVariable {
										variablesSet[_var.Name.Value] = struct{}{}
									}
								}
							}
						}
					}
				}
			}

			// If the field has a selection set, extract its variables too.
			if field.SelectionSet != nil {
				extractUsedVariablesNames(*field.SelectionSet, variablesSet, fragmentsPool)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			extractUsedVariablesNames(*inlineFrag.GetSelections(), variablesSet, fragmentsPool)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			extractUsedVariablesNames(fragmentsPool[fragSpread.FragmentName.Value].SelectionSet,
				variablesSet,
				fragmentsPool)
		}
	}
}

func getFragmentsPool(doc document) map[string]*fragmentDefinition {
	fragmentsPool := make(map[string]*fragmentDefinition)

	// Populate the fragment dictionary with all the fragments in the document.
	// the key of the dictionary is the name of the fragment definition for easy
	// access.
	for _, def := range doc.Definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			fragmentsPool[fragDef.FragmentName.Value] = fragDef
		}
	}

	return fragmentsPool
}

func detectFragmentCycles(fragDef fragmentDefinition,
	visited map[string]struct{},
	fragmentsPool map[string]*fragmentDefinition) {
	// spreads is a set that contains all fragment spreads descendants of fragDef.
	spreads := make(map[string]struct{}, 0)

	// Extract all fragment spreads descendants of fragDef.
	extractUsedFragmentsNames(fragDef.GetSelectionSet(), spreads)

	// For each spread, make sure that it is not already exists in visited.
	// If it is, panic.
	for spread := range spreads {
		if _, ok := visited[spread]; ok {
			panic(errors.New("The graph of fragment spreads must not" +
				" form any cycles including spreading itself"))
		}

		// Add the spread to the visited set.
		visited[spread] = struct{}{}

		// Call detectFragmentCycles with the target of spread.
		detectFragmentCycles(*fragmentsPool[spread], visited, fragmentsPool)
	}
}

func isInputType(schema document, variableType _type) bool {
	basicScalars := map[string]struct{}{
		"Boolean": struct{}{},
		"String":  struct{}{},
		"Int":     struct{}{},
		"Float":   struct{}{},
	}

	if _, isBasicType := basicScalars[variableType.GetTypeName()]; isBasicType {
		return true
	}

	if nonNullType, isNonNullType := variableType.(*nonNullType); isNonNullType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isInputType(schema, nonNullType.OfType)
	}

	if listType, isListType := variableType.(*listType); isListType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isInputType(schema, listType.OfType)
	}

	for _, def := range schema.Definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.GetName().Value == variableType.GetTypeName() {
				switch typeDef.(type) {
				case *scalarTypeDefinition,
					*enumTypeDefinition,
					*inputObjectTypeDefinition:
					{
						return true
					}
				default:
					{
						return false
					}
				}
			}
		}
	}

	return false
}

func isOutputType(schema document, variableType _type) bool {
	basicScalars := map[string]struct{}{
		"Boolean": struct{}{},
		"String":  struct{}{},
		"Int":     struct{}{},
		"Float":   struct{}{},
	}

	if _, isBasicType := basicScalars[variableType.GetTypeName()]; isBasicType {
		return true
	}

	if nonNullType, isNonNullType := variableType.(*nonNullType); isNonNullType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isOutputType(schema, nonNullType.OfType)
	}

	if listType, isListType := variableType.(*listType); isListType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isOutputType(schema, listType.OfType)
	}

	for _, def := range schema.Definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.GetName().Value == variableType.GetTypeName() {
				switch typeDef.(type) {
				case *scalarTypeDefinition,
					*objectTypeDefinition,
					*interfaceTypeDefinition,
					*unionTypeDefinition,
					*enumTypeDefinition:
					{
						return true
					}
				default:
					{
						return false
					}
				}
			}
		}
	}

	return false
}
