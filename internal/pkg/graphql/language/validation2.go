package language

import (
	"fmt"
	"github.com/pkg/errors"
)

type directiveUsage struct {
	directive *directive
	location  executableDirectiveLocation
}

type variableUsage struct {
	selectionSet selectionSet
	field        *field
	directive    *directive
	argument     argument
	objectField  *objectField
	listValue    *listValue
	variable     variable
}

type objectValueUsage struct {
	locationType _type
	object       objectValue
}

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
func validateFragmentSpreadIsPossible(schema, doc document) {
	fragmentsPool := getFragmentsPool(doc)
	rootQueryTypeDef := getRootQueryTypeDefinition(schema)

	for _, def := range doc.Definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			checkSpreadsPossibilityInSelectionSet(
				schema,
				exeDef.GetSelectionSet(),
				exeDef.GetSelectionSet(),
				rootQueryTypeDef,
				fragmentsPool,
			)
		}
	}
}

func checkSpreadsPossibilityInSelectionSet(
	schema document,
	rootSelectionSet selectionSet,
	selectionSet selectionSet,
	parentType typeDefinition,
	fragmentsPool map[string]*fragmentDefinition,
) {
	// TODO: implement
}

func getPossibleTypes(schema document, t _type) map[string]struct{} {
	typesSet := make(map[string]struct{})
	typeDef := getTypeDefinitionByType(schema, t)

	switch v := typeDef.(type) {
	case *objectTypeDefinition:
		typesSet[v.Name.Value] = struct{}{}
	case *interfaceTypeDefinition:
		for _, def := range schema.Definitions {
			if objectTypeDef, isObjectTypeDef := def.(*objectTypeDefinition); isObjectTypeDef {
				if objectTypeDef.ImplementsInterfaces != nil {
					for _, iface := range *objectTypeDef.ImplementsInterfaces {
						if iface.GetTypeName() == v.Name.Value {
							typesSet[v.Name.Value] = struct{}{}
						}
					}
				}
			}
		}
	case *unionTypeDefinition:
		if v.UnionMemberTypes != nil {
			for _, unionMember := range *v.UnionMemberTypes {
				typesSet[unionMember.GetTypeName()] = struct{}{}
			}
		}
	default:
		panic(errors.New("Cannot get possible types of a type which is not " +
			"an object, interface or union type definition"))
	}

	return typesSet
}

// http://spec.graphql.org/draft/#sec-Values
func validateValuesOfCorrectType(doc document) {
	// TODO: Implement
}

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Names
func validateInputObjectFieldNames(schema document, doc document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.Definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			objectValueUsages := collectObjectValueUsages(
				schema,
				exeDef.GetSelectionSet(),
				exeDef.GetSelectionSet(),
				fragmentsPool,
			)

			for _, usage := range objectValueUsages {
				typeDef := getTypeDefinitionByType(schema, usage.locationType)

				inputObjectTypeDef, isInputObjectTypeDef := typeDef.(*inputObjectTypeDefinition)

				if isInputObjectTypeDef {
					if inputObjectTypeDef.InputFieldsDefinition != nil {
						for _, inputField := range usage.object.Values {
							if !isInputFieldDefined(
								inputField,
								*inputObjectTypeDef.InputFieldsDefinition,
							) {
								panic(errors.New("Every input field provided in an input" +
									" object value must be defined in the set of possible fields of" +
									" that input object’s expected type."))
							}
						}
					} else {
						panic(errors.New("Cannot find input field definition in empty input" +
							" object type definition"))
					}
				}
			}
		}
	}
}

func isInputFieldDefined(field objectField, fieldsDefinition inputFieldsDefinition) bool {
	for _, fieldDef := range fieldsDefinition {
		if fieldDef.Name.Value == field.Name.Value {
			return true
		}
	}

	return false
}

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Uniqueness
func validateInputObjectFieldUniqueness(doc document) {
	for _, def := range doc.Definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			inputObjects := collectInputObjects(exeDef.GetSelectionSet())

			for _, inputObject := range inputObjects {
				fields := make(map[string]struct{})

				for _, inputField := range inputObject.Values {
					if _, isFieldAlreadyExists := fields[inputField.Name.Value]; isFieldAlreadyExists {
						panic(errors.New("Input objects must not contain" +
							" more than one field of the same name, otherwise an" +
							" ambiguity would exist which includes an ignored" +
							" portion of syntax"))
					}

					fields[inputField.Name.Value] = struct{}{}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Input-Object-Required-Fields
func validateInputObjectRequiredFields(schema document, doc document) {
	fragmentsPool := getFragmentsPool(doc)

	// For each definition in the document:
	for _, def := range doc.Definitions {
		// If it is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			// Collect input object values from the executable definition.
			inputObjectUsages := collectObjectValueUsages(
				schema,
				exeDef.GetSelectionSet(),
				exeDef.GetSelectionSet(),
				fragmentsPool,
			)

			// For each input object usage:
			for _, inputObjectUsage := range inputObjectUsages {
				// Get the type definition that defines the input object.
				typeDef := getTypeDefinitionByType(schema, inputObjectUsage.locationType)

				// If the type definition is an input object type definition:
				if inputObjectTypeDef, isInputObjectTypeDef := typeDef.(*inputObjectTypeDefinition); isInputObjectTypeDef {
					// If the input object type definition has input fields definition:
					if inputObjectTypeDef.InputFieldsDefinition != nil {
						// For each field definition:
						for _, fieldDef := range *inputObjectTypeDef.InputFieldsDefinition {
							_, isNonNull := fieldDef.Type.(*nonNullType)

							// If the field's expected type is a non null type and it does not
							// have a default value, the field must exist in the object usage:
							if isNonNull && fieldDef.DefaultValue == nil {
								// Set a flag that indicates if the required field has been found
								// in the object usage
								requiredFieldExists := false

								// For each field in the object usage:
								for _, inputField := range inputObjectUsage.object.Values {
									// If the field's name equals to the field definition's name:
									if inputField.Name.Value == fieldDef.Name.Value {
										requiredFieldExists = true

										_, isNull := inputField.Value.(*nullValue)

										// If the value is null, panic
										if isNull {
											panic(errors.New("An input field is required " +
												"if it has a non‐null type and does not have a default" +
												" value."))
										}
									}
								}

								// If the required field does not exist in the object usage, panic.
								if !requiredFieldExists {
									panic(errors.New("All required fields in input object must exist"))
								}
							}
						}
					}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Directives-Are-Defined
func validateDirectivesAreDefined(schema, doc document) {
	// For each definition in the document,
	for _, def := range doc.Definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			if exeDef.GetDirectives() != nil {
				// For each directive in the executable definition, check if
				// it is defined.
				for _, directive := range *exeDef.GetDirectives() {
					isDirectiveDefined(schema, directive)
				}
			}

			// If the definition is an operation definition, it is possible to
			// use directives in the context of the operation's variables.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				if opDef.VariableDefinitions != nil {
					// For each variable definition, check if its directives
					// are defined.
					for _, varDef := range *opDef.VariableDefinitions {
						if varDef.Directives != nil {
							for _, directive := range *varDef.Directives {
								isDirectiveDefined(schema, directive)
							}
						}
					}
				}
			}

			// Check (recursively) the directives in the operation's selection set.
			checkIfDirectivesInSelectionSetAreDefined(schema, exeDef.GetSelectionSet())
		}
	}
}

// http://spec.graphql.org/draft/#sec-Directives-Are-In-Valid-Locations
func validateDirectivesAreInValidLocations(schema document, doc document) {
	// For each definition in the document:
	for _, def := range doc.Definitions {
		// Create a list of directive usages
		directiveUsages := make([]*directiveUsage, 0)

		// If the definition is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			var dirLocation executableDirectiveLocation

			// If the executable definition is an operation definition,
			// set dirLocation to be the type of the operation.
			// Else, set dirLocation to be fragment definition.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				switch opDef.OperationType {
				case operationQuery:
					{
						dirLocation = edlQuery
					}
				case operationMutation:
					{
						dirLocation = edlMutation
					}
				case operationSubscription:
					{
						dirLocation = edlSubscription
					}
				}

				// If the operation definition has variable definitions:
				if opDef.VariableDefinitions != nil {
					// For each variable definition:
					for _, varDef := range *opDef.VariableDefinitions {
						// If the variable definition has directives:
						if varDef.Directives != nil {
							// For each directive:
							for _, dir := range *varDef.Directives {
								// If the directive used in a location that it is not
								// defined to be used, panic.
								if !checkDirectiveLocation(schema,
									&directiveUsage{&dir, edlVariableDefinition}) {
									panic(errors.New("GraphQL servers define what directives " +
										"they support and where they support them. For each usage of a " +
										"directive, the directive must be used in a location that the " +
										"server has declared support for."))
								}
							}
						}
					}
				}
			} else if _, isFragDef := exeDef.(*fragmentDefinition); isFragDef {
				dirLocation = edlFragmentDefinition
			}

			// If the definition has directives:
			if exeDef.GetDirectives() != nil {
				// For each directive:
				for _, dir := range *exeDef.GetDirectives() {
					// If the directive used in a location that it is not
					// defined to be used, panic.
					if !checkDirectiveLocation(schema, &directiveUsage{&dir, dirLocation}) {
						panic(errors.New("GraphQL servers define what directives " +
							"they support and where they support them. For each usage of a " +
							"directive, the directive must be used in a location that the " +
							"server has declared support for."))
					}
				}
			}

			// Get a list of all directive usages in the executable definition's selection set.
			directiveUsages = append(directiveUsages,
				extractDirectivesWithLocationsFromSelectionSet(exeDef.GetSelectionSet())...)

			// For each directive usage:
			for _, usage := range directiveUsages {
				// If the directive used in a location that it is not
				// defined to be used, panic.
				if !checkDirectiveLocation(schema, usage) {
					panic(errors.New("GraphQL servers define what directives " +
						"they support and where they support them. For each usage of a " +
						"directive, the directive must be used in a location that the " +
						"server has declared support for."))
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Directives-Are-Unique-Per-Location
func validateDirectivesAreUniquePerLocation(doc document) {
	errMsg := "Directives are used to describe" +
		" some metadata or behavioral change on the definition" +
		" they apply to. When more than one directive of the same" +
		" name is used, the expected metadata or behavior becomes" +
		" ambiguous, therefore only one of each directive is " +
		"allowed per location."

	// For each definition in the document:
	for _, def := range doc.Definitions {
		// If the definition is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {

			// If the executable definition is an operation definition,
			// set dirLocation to be the type of the operation.
			// Else, set dirLocation to be fragment definition.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				// If the operation definition has variable definitions:
				if opDef.VariableDefinitions != nil {
					// For each variable definition:
					for _, varDef := range *opDef.VariableDefinitions {
						// If the variable definition has directives:
						if varDef.Directives != nil {
							if !checkDirectivesUniqueness(*varDef.Directives) {
								panic(errors.New(errMsg))
							}
						}
					}
				}
			}

			// If the definition has directives:
			if exeDef.GetDirectives() != nil {
				if !checkDirectivesUniqueness(*exeDef.GetDirectives()) {
					panic(errors.New(errMsg))
				}
			}

			if !checkDirectivesUniquenessInSelectionSet(exeDef.GetSelectionSet()) {
				panic(errors.New(errMsg))
			}
		}
	}
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
			variableUsages := make(map[interface{}]*variableUsage)
			fragmentsPool := getFragmentsPool(doc)
			collectVariableUsages(opDef.SelectionSet, variableUsages, fragmentsPool)

			for _, varUsage := range variableUsages {
				for _, varDef := range *opDef.VariableDefinitions {
					if varUsage.variable.Name.Value == varDef.Variable.Name.Value {
						if !isVariableUsageAllowed(schema, &varDef, varUsage, opDef.SelectionSet, fragmentsPool) {
							panic(errors.New("Variable usages must be compatible" +
								" with the arguments they are passed to"))
						}
					}
				}
			}
		}
	}
}

func isVariableUsageAllowed(
	schema document,
	varDef *variableDefinition,
	variableUsage *variableUsage,
	rootSelectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) bool {
	var (
		hasNonNullVariableDefaultValue bool
		hasLocationDefaultValue        bool
	)

	var argDef *inputValueDefinition

	// Get the argument definition from the schema according to the variable
	// usage.
	if variableUsage.directive != nil {
		argDef = getDirectiveArgumentDefinition(schema, *variableUsage.directive, variableUsage.argument)
	} else {
		argDef = getFieldArgumentDefinition(
			schema,
			rootSelectionSet,
			*variableUsage.field,
			variableUsage.argument,
			fragmentsPool,
		)
	}

	// Let locationType be the expected type of the Argument, ObjectField, or ListValue
	// entry where variableUsage is located.
	locationType := argDef.Type

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
		if argDef.DefaultValue != nil {
			hasLocationDefaultValue = true
		}

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

func getDirectiveArgumentDefinition(schema document, directive directive, argument argument) *inputValueDefinition {
	for _, def := range schema.Definitions {
		if directiveDef, isDirectiveDef := def.(*directiveDefinition); isDirectiveDef {
			if directiveDef.ArgumentsDefinition != nil {
				for _, argDef := range *directiveDef.ArgumentsDefinition {
					if argDef.Name.Value == argument.Name.Value {
						return &argDef
					}
				}
			}
		}
	}

	panic(errors.New("could not find argument in any directive definition"))
}

func getFieldArgumentDefinition(
	schema document,
	rootSelectionSet selectionSet,
	field field,
	argument argument,
	fragmentsPool map[string]*fragmentDefinition,
) *inputValueDefinition {
	selectionFieldDef := getFieldDefinitionByFieldSelection(
		getRootQueryTypeDefinition(schema),
		&field,
		rootSelectionSet,
		schema,
		fragmentsPool,
	)

	if selectionFieldDef.ArgumentsDefinition != nil {
		for _, arg := range *selectionFieldDef.ArgumentsDefinition {
			if arg.Name.Value == argument.Name.Value {
				return &arg
			}
		}

		panic(errors.New(fmt.Sprintf("could not find argument %s in field definition named %s",
			argument.Name.Value,
			selectionFieldDef.Name.Value)))
	} else {
		panic(errors.New("cannot find an argument definition in a field definition that " +
			"does not contain argument definitions"))
	}
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

	// Otherwise, if variableType is a list type, return false.
	if isVariableTypeAListType {
		return false
	}

	// Return true if variableType and locationType are identical, otherwise false.
	return variableType.GetTypeName() == locationType.GetTypeName()
}

func collectObjectValueUsages(
	schema document,
	rootSelectionSet selectionSet,
	selectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) []objectValueUsage {
	usages := make([]objectValueUsage, 0)

	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection has directives, check their arguments too.
		if selection.GetDirectives() != nil {
			for _, directive := range *selection.GetDirectives() {
				if directive.Arguments != nil {
					// For each argument, if it is a variable, add it to the variables set.
					// If it is a list, check each item in the list.
					// If it is an object, check each value in the object.
					for _, arg := range *directive.Arguments {
						if object, isObject := arg.Value.(*objectValue); isObject {
							inputValueDef := getDirectiveArgumentDefinition(schema, directive, arg)

							usages = append(usages, objectValueUsage{
								locationType: inputValueDef.Type,
								object:       *object,
							})
						}
					}
				}
			}
		}

		// If the selection is a field, extract its variable usages
		// (the field's arguments, and the field's directive's arguments).
		if field, isField := selection.(*field); isField {
			if field.Arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.Arguments {
					if object, isObject := arg.Value.(*objectValue); isObject {
						inputValueDef := getFieldArgumentDefinition(
							schema,
							rootSelectionSet,
							*field,
							arg,
							fragmentsPool,
						)

						usages = append(usages, objectValueUsage{
							locationType: inputValueDef.Type,
							object:       *object,
						})
					}
				}
			}

			// If the field has a selection set, extract its variables too.
			if field.SelectionSet != nil {
				selectionSetObjectUsages := collectObjectValueUsages(
					schema,
					rootSelectionSet,
					*field.SelectionSet,
					fragmentsPool,
				)

				usages = append(usages, selectionSetObjectUsages...)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			inlineFragObjectUsages := collectObjectValueUsages(
				schema,
				rootSelectionSet,
				inlineFrag.SelectionSet,
				fragmentsPool,
			)

			usages = append(usages, inlineFragObjectUsages...)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			inlineFragObjectUsages := collectObjectValueUsages(
				schema,
				rootSelectionSet,
				fragmentsPool[fragSpread.FragmentName.Value].SelectionSet,
				fragmentsPool,
			)

			usages = append(usages, inlineFragObjectUsages...)
		}
	}

	return usages
}

func collectVariableUsages(
	selectionSet selectionSet,
	usages map[interface{}]*variableUsage,
	fragmentsPool map[string]*fragmentDefinition,
) {
	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection has directives, check their arguments too.
		if selection.GetDirectives() != nil {
			for _, directive := range *selection.GetDirectives() {
				if directive.Arguments != nil {
					// For each argument, if it is a variable, add it to the variables set.
					// If it is a list, check each item in the list.
					// If it is an object, check each value in the object.
					for _, arg := range *directive.Arguments {
						if _var, isVariable := arg.Value.(*variable); isVariable {
							usages[&arg] = &variableUsage{
								field:       nil,
								directive:   &directive,
								argument:    arg,
								listValue:   nil,
								objectField: nil,
								variable:    *_var,
							}
						} else if list, isList := arg.Value.(*listValue); isList {
							for _, item := range list.Values {
								if _var, isVariable := item.(*variable); isVariable {
									usages[&list] = &variableUsage{
										field:       nil,
										directive:   &directive,
										argument:    arg,
										listValue:   list,
										objectField: nil,
										variable:    *_var,
									}
								}
							}
						} else if object, isObject := arg.Value.(*objectValue); isObject {
							for _, field := range object.Values {
								if _var, isVariable := field.Value.(*variable); isVariable {
									usages[&object] = &variableUsage{
										field:       nil,
										directive:   &directive,
										argument:    arg,
										listValue:   nil,
										objectField: &field,
										variable:    *_var,
									}
								}
							}
						}
					}
				}
			}
		}

		// If the selection is a field, extract its variable usages
		// (the field's arguments, and the field's directive's arguments).
		if field, isField := selection.(*field); isField {
			if field.Arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.Arguments {
					if _var, isVariable := arg.Value.(*variable); isVariable {
						usages[&arg] = &variableUsage{
							field:       field,
							directive:   nil,
							argument:    arg,
							listValue:   nil,
							objectField: nil,
							variable:    *_var,
						}
					} else if list, isList := arg.Value.(*listValue); isList {
						for _, item := range list.Values {
							if _var, isVariable := item.(*variable); isVariable {
								usages[&list] = &variableUsage{
									field:       field,
									directive:   nil,
									argument:    arg,
									listValue:   list,
									objectField: nil,
									variable:    *_var,
								}
							}
						}
					} else if object, isObject := arg.Value.(*objectValue); isObject {
						for _, objectField := range object.Values {
							if _var, isVariable := objectField.Value.(*variable); isVariable {
								usages[&object] = &variableUsage{
									field:       field,
									directive:   nil,
									argument:    arg,
									listValue:   nil,
									objectField: &objectField,
									variable:    *_var,
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
		// If the selection has directives, check their arguments too.
		if selection.GetDirectives() != nil {
			for _, directive := range *selection.GetDirectives() {
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

func detectFragmentCycles(
	fragDef fragmentDefinition,
	visited map[string]struct{},
	fragmentsPool map[string]*fragmentDefinition,
) {
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

func checkIfDirectivesInSelectionSetAreDefined(schema document, selectionSet selectionSet) {
	// For each selection, look for directives.
	for _, selection := range selectionSet {
		// And the field contains directives
		if selection.GetDirectives() != nil {
			// For each directive, decide whether it is defined or not.
			for _, directive := range *selection.GetDirectives() {
				// If the directive is not defined (which means the directive is not a built in
				// directive and we could not find a proper directive definition in
				// the schema), return false.
				isDirectiveDefined(schema, directive)
			}
		}

		// If the field contains a selection set, check it's directives too.
		if selection.GetSelections() != nil {
			checkIfDirectivesInSelectionSetAreDefined(schema, *selection.GetSelections())
		}
	}
}

func isDirectiveDefined(schema document, directive directive) {
	// A map that contains all graphql's built in directives for quick access.
	builtInDirectiveNames := map[string]struct{}{
		"skip":       struct{}{},
		"include":    struct{}{},
		"deprecated": struct{}{},
	}

	// Check if the directive is a build in directive. If it is, turn on the
	// flag. Else, Search the directive in the schema.
	if _, isBuiltInDirective := builtInDirectiveNames[directive.Name.Value]; isBuiltInDirective {
		return
	} else {
		// For each definition in the schema, if it is a directive definition,
		// compare its name to the current directive from the document.
		for _, def := range schema.Definitions {
			if directiveDef, isDirectiveDef := def.(*directiveDefinition); isDirectiveDef {
				// If the name are equal, turn on the flag.
				if directiveDef.Name.Value == directive.Name.Value {
					return
				}
			}
		}
	}

	panic(errors.New("GraphQL servers define what directives they " +
		"support. For each usage of a directive, the directive must be " +
		"available on that server"))
}

func collectInputObjects(selectionSet selectionSet) []*objectValue {
	inputObjects := make([]*objectValue, 0)

	for _, selection := range selectionSet {
		if field, isField := selection.(*field); isField {
			if field.Arguments != nil {
				for _, arg := range *field.Arguments {
					if object, isObject := arg.Value.(*objectValue); isObject {
						inputObjects = append(inputObjects, object)
					}
				}
			}
		}

		if selection.GetDirectives() != nil {
			for _, directive := range *selection.GetDirectives() {
				if directive.Arguments != nil {
					for _, arg := range *directive.Arguments {
						if object, isObject := arg.Value.(*objectValue); isObject {
							inputObjects = append(inputObjects, object)
						}
					}
				}
			}
		}

		if selection.GetSelections() != nil {
			inputObjects = append(inputObjects, collectInputObjects(*selection.GetSelections())...)
		}
	}

	return inputObjects
}

func checkDirectiveLocation(schema document, usage *directiveUsage) bool {
	// If the directive is one of the built in directives, and it has been used
	// on a field, fragment spread or an inline fragment return true.
	if usage.directive.Name.Value == "skip" ||
		usage.directive.Name.Value == "include" {
		if usage.location == edlField || usage.location == edlFragmentSpread ||
			usage.location == edlInlineFragment {
			return true
		}
	}

	//For each definition in the schema:
	for _, def := range schema.Definitions {
		// If the definition is a directive definition:
		if dirDef, isDirDef := def.(*directiveDefinition); isDirDef {
			// If the name of the directive definition equals to the name of the
			// given directive:
			if usage.directive.Name.Value == dirDef.Name.Value {
				// For each location in the directive definition, compare it with the
				// given location. If they are equal return true.
				for _, location := range dirDef.DirectiveLocations {
					// TODO: change directiveLocation type from string to
					// TODO: executableDirectiveLocation in order to remove the
					// TODO: conversion to string.
					if string(usage.location) == string(location) {
						return true
					}
				}

				// If we finished the directive locations loop and we did not
				// return true, it is useless to keep searching for another
				// directive definitions because directive names are unique.
				break
			}
		}
	}

	// If we arrived here it means that we could not find any directive location
	// that is equal to the given usage location.
	return false
}

func extractDirectivesWithLocationsFromSelectionSet(selectionSet selectionSet) []*directiveUsage {
	usages := make([]*directiveUsage, 0)

	// For each selection in the selection set:
	for _, selection := range selectionSet {
		// If the selection has directives:
		if selection.GetDirectives() != nil {
			// For each directive, set the location to be the kind of the selection.
			for _, dir := range *selection.GetDirectives() {
				var location executableDirectiveLocation

				switch selection.(type) {
				case *field:
					{
						location = edlField
					}
				case *fragmentSpread:
					{
						location = edlFragmentSpread
					}
				case *inlineFragment:
					{
						location = edlInlineFragment
					}
				}

				// Append the directive usage to the usages list.
				usages = append(usages, &directiveUsage{&dir, location})
			}
		}

		// If the selection has a selection set, get a list of all of its directive usages.
		if selection.GetSelections() != nil {
			usages = append(usages, extractDirectivesWithLocationsFromSelectionSet(*selection.GetSelections())...)
		}
	}

	// Return the usages list.
	return usages
}

func checkDirectivesUniqueness(directives directives) bool {
	directivesSet := make(map[string]struct{})

	// For each directive in directives.
	for _, dir := range directives {
		// If the directive name already exists in the unique set, return false.
		if _, isDirNameAlreadyExists := directivesSet[dir.Name.Value]; isDirNameAlreadyExists {
			return false
		}

		// Insert the directive name to the unique set.
		directivesSet[dir.Name.Value] = struct{}{}
	}

	return true
}

func checkDirectivesUniquenessInSelectionSet(selectionSet selectionSet) bool {
	// For each selection in the selection set:
	for _, selection := range selectionSet {
		// If the selection has directives, check their uniqueness.
		if selection.GetDirectives() != nil {
			if !checkDirectivesUniqueness(*selection.GetDirectives()) {
				return false
			}
		}

		if selection.GetSelections() != nil {
			if !checkDirectivesUniquenessInSelectionSet(*selection.GetSelections()) {
				return false
			}
		}
	}

	return true
}

func getRootQueryTypeDefinition(schema document) *objectTypeDefinition {
	rootQueryTypeName := "Query"

	for _, def := range schema.Definitions {
		if schemaDef, isSchemaDef := def.(*schemaDefinition); isSchemaDef {
			for _, rootOperationType := range schemaDef.RootOperationTypeDefinitions {
				if rootOperationType.OperationType == operationQuery {
					rootQueryTypeName = rootOperationType.NamedType.Value
				}
			}
		}
	}

	for _, def := range schema.Definitions {
		if objectTypeDef, isObjectTypeDef := def.(*objectTypeDefinition); isObjectTypeDef {
			if objectTypeDef.Name.Value == rootQueryTypeName {
				return objectTypeDef
			}
		}
	}

	panic(errors.New("could not find root query type"))
}

func getSelectionSetType(
	parentType typeDefinition,
	target *selectionSet,
	current selectionSet,
	schema document,
	fragmentsPool map[string]*fragmentDefinition,
) _type {
	var typeToReturn _type
	for _, selection := range current {
		switch s := selection.(type) {
		case *field:
			if s.SelectionSet != nil {
				var fieldType _type
				switch t := parentType.(type) {
				case *objectTypeDefinition:
					if t.FieldsDefinition != nil {
						for _, fieldDef := range *t.FieldsDefinition {
							if fieldDef.Name.Value == s.Name.Value {
								fieldType = fieldDef.Type
								break
							}
						}
					}
				case *interfaceTypeDefinition:
					if t.FieldsDefinition != nil {
						for _, fieldDef := range *t.FieldsDefinition {
							if fieldDef.Name.Value == s.Name.Value {
								fieldType = fieldDef.Type
								break
							}
						}
					}
				case *unionTypeDefinition:
					if t.UnionMemberTypes != nil {
						for _, unionMember := range *t.UnionMemberTypes {
							fieldType = getSelectionSetType(
								getTypeDefinitionByType(schema, &unionMember),
								target,
								*s.SelectionSet,
								schema,
								fragmentsPool)

							if fieldType != nil {
								break
							}
						}
					}
				}

				if s.SelectionSet == target {
					return fieldType
				} else {
					fieldSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, fieldType),
						target,
						*s.SelectionSet,
						schema,
						fragmentsPool,
					)

					if fieldSelectionSetType != nil {
						typeToReturn = fieldSelectionSetType
					}
				}
			} else {
				return nil
			}
		case *inlineFragment:
			if s.SelectionSet != nil {
				if &s.SelectionSet == target {
					return &s.TypeCondition.NamedType
				} else {
					inlineFragSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, &s.TypeCondition.NamedType),
						target,
						s.SelectionSet,
						schema,
						fragmentsPool,
					)

					if inlineFragSelectionSetType != nil {
						typeToReturn = inlineFragSelectionSetType
					}
				}
			}
		case *fragmentSpread:
			frag := fragmentsPool[s.FragmentName.Value]

			if frag.SelectionSet != nil {
				if &frag.SelectionSet == target {
					return &frag.TypeCondition.NamedType
				} else {
					fragSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, &frag.TypeCondition.NamedType),
						target,
						frag.SelectionSet,
						schema,
						fragmentsPool,
					)

					if fragSelectionSetType != nil {
						typeToReturn = fragSelectionSetType
					}
				}
			}
		}
	}

	if typeToReturn != nil {
		return typeToReturn
	} else {
		panic(errors.New("could not find the requested selection set type."))
	}
}

func getFieldDefinitionByFieldSelection(
	parentType typeDefinition,
	targetSelection selection,
	selectionSet selectionSet,
	schema document,
	fragmentsPool map[string]*fragmentDefinition,
) *fieldDefinition {
	var tachlessFieldDefinition *fieldDefinition

	for _, selection := range selectionSet {
		switch s := selection.(type) {
		case *field:
			switch t := parentType.(type) {
			case *objectTypeDefinition:
				if t.FieldsDefinition != nil {
					for _, fieldDef := range *t.FieldsDefinition {
						if fieldDef.Name.Value == s.Name.Value {
							tachlessFieldDefinition = &fieldDef
							break
						}
					}
				}
			case *interfaceTypeDefinition:
				if t.FieldsDefinition != nil {
					for _, fieldDef := range *t.FieldsDefinition {
						if fieldDef.Name.Value == s.Name.Value {
							tachlessFieldDefinition = &fieldDef
							break
						}
					}
				}
			case *unionTypeDefinition:
				if t.UnionMemberTypes != nil {
					for _, unionMember := range *t.UnionMemberTypes {
						fieldDef := getFieldDefinitionByFieldSelection(
							getTypeDefinitionByType(schema, &unionMember),
							targetSelection,
							selectionSet,
							schema,
							fragmentsPool)

						if fieldDef != nil {
							return fieldDef
						}
					}

					return nil
				}
			}

			if tachlessFieldDefinition == nil {
				panic(errors.New("could not find a field definition named " + s.Name.Value))
			}

			if selection == targetSelection {
				return tachlessFieldDefinition
			} else if s.GetSelections() != nil {
				return getFieldDefinitionByFieldSelection(
					getTypeDefinitionByType(schema, tachlessFieldDefinition.Type),
					targetSelection,
					*s.GetSelections(),
					schema,
					fragmentsPool,
				)
			} else {
				return nil
			}
		case *inlineFragment:
			return getFieldDefinitionByFieldSelection(
				getTypeDefinitionByType(schema, &s.TypeCondition.NamedType),
				targetSelection,
				*s.GetSelections(),
				schema,
				fragmentsPool,
			)
		case *fragmentSpread:
			return getFieldDefinitionByFieldSelection(
				getTypeDefinitionByType(
					schema,
					&fragmentsPool[s.FragmentName.Value].TypeCondition.NamedType,
				),
				targetSelection,
				*s.GetSelections(),
				schema, fragmentsPool,
			)
		}
	}

	panic(errors.New("empty selection cannot query for selection type"))
}

func getTypeDefinitionByType(schema document, t _type) typeDefinition {
	for _, def := range schema.Definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.GetName().Value == t.GetTypeName() {
				return typeDef
			}
		}
	}

	panic(errors.New("could not find a type definition named: " + t.GetTypeName()))
}
