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
	argument     *argument
	objectField  *objectField
	listValue    *listValue
	variable     variable
}

type objectValueUsage struct {
	locationType _type
	object       objectValue
}

/**************************/
/** Validation Functions **/
/**************************/

// http://spec.graphql.org/draft/#sec-Fragments-Must-Be-Used
func validateFragmentsMustBeUsed(doc *document) {
	fragmentSpreadTargets := make(map[string]struct{}, 0)

	// Extract fragment spread targets from all the definitions in the document.
	for _, def := range doc.definitions {
		exeDef := def.(executableDefinition)
		extractUsedFragmentsNames(exeDef.SelectionSet(), fragmentSpreadTargets)
	}

	// For each fragment definition in the document, check if it's name exists
	// in the fragment spread targets set. If not, panic.
	for _, def := range doc.definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			if _, ok := fragmentSpreadTargets[fragDef.fragmentName.value]; !ok {
				panic(errors.New("Defined fragments must be used within " +
					"a document: \"" + fragDef.fragmentName.value + "\" is not used"))
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Fragment-spread-target-defined
func validateFragmentSpreadTargetDefined(doc *document) {
	fragmentSpreadTargets := make(map[string]struct{}, 0)
	fragmentDefinitionsNames := make(map[string]struct{})

	// Collect all the fragment definitions' names into a dictionary of names.
	for _, def := range doc.definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			fragmentDefinitionsNames[fragDef.fragmentName.value] = struct{}{}
		}
	}

	// Extract fragment spread targets from all the definitions in the document.
	for _, def := range doc.definitions {
		exeDef := def.(executableDefinition)
		extractUsedFragmentsNames(exeDef.SelectionSet(), fragmentSpreadTargets)
	}

	// Verify that all the targets in the targets set are referring to a defined fragment.
	// If not, panic.
	for target := range fragmentSpreadTargets {
		if _, ok := fragmentDefinitionsNames[target]; !ok {
			panic(errors.New("Named fragment spreads must refer to fragments" +
				" defined within the document: \"" + target + "\" is not defined"))
		}
	}
}

// http://spec.graphql.org/draft/#sec-Fragment-spreads-must-not-form-cycles
func validateFragmentSpreadsMustNotFormCycles(doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	// For each fragment, call detectFragmentCycles()
	for _, fragDef := range fragmentsPool {
		visited := make(map[string]struct{}, 0)
		detectFragmentCycles(*fragDef, visited, fragmentsPool)
	}
}

// http://spec.graphql.org/draft/#sec-Fragment-spread-is-possible
func validateFragmentSpreadIsPossible(schema, doc *document) {
	// Get the fragments pool for quick access.
	fragmentsPool := getFragmentsPool(doc)

	// Get the root query type from the schema.
	rootQueryTypeDef := getRootQueryTypeDefinition(schema)

	// For each operation in the document, check its spreads' possibility.
	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			checkSpreadsPossibilityInSelectionSet(
				schema,
				opDef.SelectionSet(),
				opDef.SelectionSet(),
				rootQueryTypeDef,
				fragmentsPool,
			)
		}
	}
}

// http://spec.graphql.org/draft/#sec-Values
func validateValuesOfCorrectType(schema, doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			checkValuesOfCorrectTypeInOperation(
				schema,
				opDef.selectionSet,
				opDef.selectionSet,
				fragmentsPool,
			)
		}
	}
}

func checkValuesOfCorrectTypeInOperation(
	schema *document,
	root selectionSet,
	set selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) {
	for _, selection := range set {
		if selection.Directives() != nil {
			for _, dir := range *selection.Directives() {
				if dir.arguments != nil {
					for _, arg := range *dir.arguments {
						argDef := getDirectiveArgumentDefinition(schema, dir, arg)
						if !assertValueType(arg.value, argDef._type) {
							panic(errors.New("Literal values must be compatible" +
								" with the type expected in the position they are found" +
								" as per the coercion rules."))
						}
					}
				}
			}
		}

		switch s := selection.(type) {
		case *field:
			if s.arguments != nil {
				for _, arg := range *s.arguments {
					argDef := getFieldArgumentDefinition(schema, root, s, arg, fragmentsPool)
					if !assertValueType(arg.value, argDef._type) {
						panic(errors.New("Literal values must be compatible" +
							" with the type expected in the position they are found" +
							" as per the coercion rules."))
					}
				}
			}
		case *inlineFragment:
			checkValuesOfCorrectTypeInOperation(
				schema,
				root,
				s.selectionSet,
				fragmentsPool,
			)
		case *fragmentSpread:
			checkValuesOfCorrectTypeInOperation(
				schema,
				root,
				fragmentsPool[s.fragmentName.value].selectionSet,
				fragmentsPool,
			)
		}
	}
}

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Names
func validateInputObjectFieldNames(schema, doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			objectValueUsages := collectObjectValueUsages(
				schema,
				exeDef.SelectionSet(),
				exeDef.SelectionSet(),
				fragmentsPool,
			)

			for _, usage := range objectValueUsages {
				typeDef := getTypeDefinitionByType(schema, usage.locationType)

				inputObjectTypeDef, isInputObjectTypeDef := typeDef.(*inputObjectTypeDefinition)

				if isInputObjectTypeDef {
					if inputObjectTypeDef.inputFieldsDefinition != nil {
						for _, inputField := range usage.object.values {
							if !isInputFieldDefined(
								inputField,
								*inputObjectTypeDef.inputFieldsDefinition,
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

// http://spec.graphql.org/draft/#sec-Input-Object-Field-Uniqueness
func validateInputObjectFieldUniqueness(doc *document) {
	for _, def := range doc.definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			inputObjects := collectInputObjects(exeDef.SelectionSet())

			for _, inputObject := range inputObjects {
				fields := make(map[string]struct{})

				for _, inputField := range inputObject.values {
					if _, isFieldAlreadyExists := fields[inputField.name.value]; isFieldAlreadyExists {
						panic(errors.New("Input objects must not contain" +
							" more than one field of the same name, otherwise an" +
							" ambiguity would exist which includes an ignored" +
							" portion of syntax"))
					}

					fields[inputField.name.value] = struct{}{}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Input-Object-Required-Fields
func validateInputObjectRequiredFields(schema, doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	// For each definition in the document:
	for _, def := range doc.definitions {
		// If it is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			// Collect input object values from the executable definition.
			inputObjectUsages := collectObjectValueUsages(
				schema,
				exeDef.SelectionSet(),
				exeDef.SelectionSet(),
				fragmentsPool,
			)

			// For each input object usage:
			for _, inputObjectUsage := range inputObjectUsages {
				// Get the type definition that defines the input object.
				typeDef := getTypeDefinitionByType(schema, inputObjectUsage.locationType)

				// If the type definition is an input object type definition:
				if inputObjectTypeDef, isInputObjectTypeDef := typeDef.(*inputObjectTypeDefinition); isInputObjectTypeDef {
					// If the input object type definition has input fields definition:
					if inputObjectTypeDef.inputFieldsDefinition != nil {
						// For each field definition:
						for _, fieldDef := range *inputObjectTypeDef.inputFieldsDefinition {
							_, isNonNull := fieldDef._type.(*nonNullType)

							// If the field's expected type is a non null type and it does not
							// have a default value, the field must exist in the object usage:
							if isNonNull && fieldDef.defaultValue == nil {
								// Set a flag that indicates if the required field has been found
								// in the object usage
								requiredFieldExists := false

								// For each field in the object usage:
								for _, inputField := range inputObjectUsage.object.values {
									// If the field's name equals to the field definition's name:
									if inputField.name.value == fieldDef.name.value {
										requiredFieldExists = true

										_, isNull := inputField._value.(*nullValue)

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
func validateDirectivesAreDefined(schema, doc *document) {
	// For each definition in the document,
	for _, def := range doc.definitions {
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			if exeDef.Directives() != nil {
				// For each directive in the executable definition, check if
				// it is defined.
				for _, directive := range *exeDef.Directives() {
					isDirectiveDefined(schema, directive)
				}
			}

			// If the definition is an operation definition, it is possible to
			// use directives in the context of the operation's variables.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				if opDef.variableDefinitions != nil {
					// For each variable definition, check if its directives
					// are defined.
					for _, varDef := range *opDef.variableDefinitions {
						if varDef.directives != nil {
							for _, directive := range *varDef.directives {
								isDirectiveDefined(schema, directive)
							}
						}
					}
				}
			}

			// Check (recursively) the directives in the operation's selection set.
			checkIfDirectivesInSelectionSetAreDefined(schema, exeDef.SelectionSet())
		}
	}
}

// http://spec.graphql.org/draft/#sec-Directives-Are-In-Valid-Locations
func validateDirectivesAreInValidLocations(schema, doc *document) {
	// For each definition in the document:
	for _, def := range doc.definitions {
		// Create a list of directive usages
		directiveUsages := make([]*directiveUsage, 0)

		// If the definition is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {
			var dirLocation executableDirectiveLocation

			// If the executable definition is an operation definition,
			// set dirLocation to be the type of the operation.
			// Else, set dirLocation to be fragment definition.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				switch opDef.operationType {
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
				if opDef.variableDefinitions != nil {
					// For each variable definition:
					for _, varDef := range *opDef.variableDefinitions {
						// If the variable definition has directives:
						if varDef.directives != nil {
							// For each directive:
							for _, dir := range *varDef.directives {
								// If the directive used in a location that it is not
								// defined to be used, panic.
								if !checkDirectiveLocation(schema,
									&directiveUsage{dir, edlVariableDefinition}) {
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
			if exeDef.Directives() != nil {
				// For each directive:
				for _, dir := range *exeDef.Directives() {
					// If the directive used in a location that it is not
					// defined to be used, panic.
					if !checkDirectiveLocation(schema, &directiveUsage{dir, dirLocation}) {
						panic(errors.New("GraphQL servers define what directives " +
							"they support and where they support them. For each usage of a " +
							"directive, the directive must be used in a location that the " +
							"server has declared support for."))
					}
				}
			}

			// Get a list of all directive usages in the executable definition's selection set.
			directiveUsages = append(directiveUsages,
				extractDirectivesWithLocationsFromSelectionSet(exeDef.SelectionSet())...)

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
func validateDirectivesAreUniquePerLocation(doc *document) {
	errMsg := "Directives are used to describe" +
		" some metadata or behavioral change on the definition" +
		" they apply to. When more than one directive of the same" +
		" name is used, the expected metadata or behavior becomes" +
		" ambiguous, therefore only one of each directive is " +
		"allowed per location."

	// For each definition in the document:
	for _, def := range doc.definitions {
		// If the definition is an executable definition:
		if exeDef, isExeDef := def.(executableDefinition); isExeDef {

			// If the executable definition is an operation definition,
			// set dirLocation to be the type of the operation.
			// Else, set dirLocation to be fragment definition.
			if opDef, isOpDef := exeDef.(*operationDefinition); isOpDef {
				// If the operation definition has variable definitions:
				if opDef.variableDefinitions != nil {
					// For each variable definition:
					for _, varDef := range *opDef.variableDefinitions {
						// If the variable definition has directives:
						if varDef.directives != nil {
							if !checkDirectivesUniqueness(*varDef.directives) {
								panic(errors.New(errMsg))
							}
						}
					}
				}
			}

			// If the definition has directives:
			if exeDef.Directives() != nil {
				if !checkDirectivesUniqueness(*exeDef.Directives()) {
					panic(errors.New(errMsg))
				}
			}

			if !checkDirectivesUniquenessInSelectionSet(exeDef.SelectionSet()) {
				panic(errors.New(errMsg))
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Variable-Uniqueness
func validateVariableUniqueness(doc *document) {
	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.variableDefinitions != nil {
				variablesSet := make(map[string]struct{}, 0)

				for _, variable := range *opDef.variableDefinitions {
					if _, isVariableExists := variablesSet[variable.variable.name.value]; isVariableExists {
						panic(errors.New("If any operation defines more than one" +
							" variable with the same name, it is ambiguous and invalid"))
					}

					variablesSet[variable.variable.name.value] = struct{}{}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-Variables-Are-Input-Types
func validateVariableAreInputTypes(schema, doc *document) {
	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.variableDefinitions != nil {
				for _, variable := range *opDef.variableDefinitions {
					if !isInputType(schema, variable._type) {
						panic(errors.New("Variables can only be input types. " +
							"Objects, unions, and interfaces cannot be used as inputs."))
					}
				}
			}
		}
	}
}

// http://spec.graphql.org/draft/#sec-All-Variable-Uses-Defined
func validateAllVariableUsesDefined(doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			usedVariablesSet := make(map[string]struct{}, 0)
			extractUsedVariablesNames(opDef.selectionSet, usedVariablesSet, fragmentsPool)

			definedVariables := make(map[string]struct{}, 0)
			if opDef.variableDefinitions != nil {
				for _, variable := range *opDef.variableDefinitions {
					definedVariables[variable.variable.name.value] = struct{}{}
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
func validateAllVariablesUsed(doc *document) {
	fragmentsPool := getFragmentsPool(doc)

	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			if opDef.variableDefinitions != nil {
				usedVariablesSet := make(map[string]struct{}, 0)
				extractUsedVariablesNames(opDef.selectionSet, usedVariablesSet, fragmentsPool)

				for _, variable := range *opDef.variableDefinitions {
					if _, isVariableUsed := usedVariablesSet[variable.variable.name.value]; !isVariableUsed {
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
func validateAllVariableUsagesAreAllowed(schema, doc *document) {
	for _, def := range doc.definitions {
		if opDef, isOpDef := def.(*operationDefinition); isOpDef {
			variableUsages := make(map[interface{}]*variableUsage)
			fragmentsPool := getFragmentsPool(doc)
			collectVariableUsages(opDef.selectionSet, variableUsages, fragmentsPool)

			for _, varUsage := range variableUsages {
				for _, varDef := range *opDef.variableDefinitions {
					if varUsage.variable.name.value == varDef.variable.name.value {
						if !isVariableUsageAllowed(schema, varDef, varUsage, opDef.selectionSet, fragmentsPool) {
							panic(errors.New("Variable usages must be compatible" +
								" with the arguments they are passed to"))
						}
					}
				}
			}
		}
	}
}

/**********************/
/** Helper Functions **/
/**********************/

func isVariableUsageAllowed(
	schema *document,
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
		argDef = getDirectiveArgumentDefinition(schema, variableUsage.directive, variableUsage.argument)
	} else {
		argDef = getFieldArgumentDefinition(
			schema,
			rootSelectionSet,
			variableUsage.field,
			variableUsage.argument,
			fragmentsPool,
		)
	}

	// Let locationType be the expected type of the Argument, ObjectField, or ListValue
	// entry where variableUsage is located.
	locationType := argDef._type

	// Let variableType be the expected type of variableDefinition
	variableType := varDef._type

	_, isVariableTypeANonNullType := variableType.(*nonNullType)
	nonNullLocationType, isLocationTypeANonNullType := locationType.(*nonNullType)

	// If locationType is a non‐null type AND variableType is NOT a non‐null type:
	if isLocationTypeANonNullType && !isVariableTypeANonNullType {
		// Let hasNonNullVariableDefaultValue be true if a default value exists for
		// variableDefinition and is not the value null.
		if varDef.defaultValue != nil {
			if _, isNullValue := (*varDef).defaultValue.value.(*nullValue); !isNullValue {
				hasNonNullVariableDefaultValue = true
			}
		}

		// Let hasLocationDefaultValue be true if a default value exists for the Argument
		// or ObjectField where variableUsage is located.
		if argDef.defaultValue != nil {
			hasLocationDefaultValue = true
		}

		// If hasNonNullVariableDefaultValue is NOT true AND hasLocationDefaultValue is
		// NOT true, return false.
		if !hasNonNullVariableDefaultValue && !hasLocationDefaultValue {
			return false
		}

		// Let nullableLocationType be the unwrapped nullable type of locationType.
		nullableLocationType := nonNullLocationType.ofType

		// Check if the types are compatible.
		return areTypesCompatible(variableType, nullableLocationType)
	}

	// Check if the types are compatible.
	return areTypesCompatible(variableType, locationType)
}

func getDirectiveArgumentDefinition(
	schema *document,
	directive *directive,
	argument *argument,
) *inputValueDefinition {
	for _, def := range schema.definitions {
		if directiveDef, isDirectiveDef := def.(*directiveDefinition); isDirectiveDef {
			if directiveDef.name.value == directive.name.value {
				if directiveDef.argumentsDefinition != nil {
					for _, argDef := range *directiveDef.argumentsDefinition {
						if argDef.name.value == argument.name.value {
							return argDef
						}
					}
				}
			}
		}
	}

	panic(errors.New("could not find argument in any directive definition"))
}

func getFieldArgumentDefinition(
	schema *document,
	rootSelectionSet selectionSet,
	field *field,
	argument *argument,
	fragmentsPool map[string]*fragmentDefinition,
) *inputValueDefinition {
	selectionFieldDef := getFieldDefinitionByFieldSelection(
		getRootQueryTypeDefinition(schema),
		field,
		rootSelectionSet,
		schema,
		fragmentsPool,
	)

	if selectionFieldDef.argumentsDefinition != nil {
		for _, arg := range *selectionFieldDef.argumentsDefinition {
			if arg.name.value == argument.name.value {
				return arg
			}
		}

		panic(errors.New(fmt.Sprintf("could not find argument %s in field definition named %s",
			argument.name.value,
			selectionFieldDef.name.value)))
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
		nullableLocationType := nonNullLocationType.ofType

		// Let nullableVariableType be the unwrapped nullable type of variableType.
		nullableVariableType := nonNullVariableType.ofType

		// Return the result of areTypesCompatible with the unwrapped types.
		return areTypesCompatible(nullableVariableType, nullableLocationType)
	}

	// Otherwise, if variableType is a non‐null type:
	if isVariableTypeIsNonNullType {
		// Let nullableVariableType be the unwrapped nullable type of variableType.
		nullableVariableType := nonNullVariableType.ofType

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
	return variableType.TypeName() == locationType.TypeName()
}

func collectObjectValueUsages(
	schema *document,
	rootSelectionSet selectionSet,
	selectionSet selectionSet,
	fragmentsPool map[string]*fragmentDefinition,
) []objectValueUsage {
	usages := make([]objectValueUsage, 0)

	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection has directives, check their arguments too.
		if selection.Directives() != nil {
			for _, directive := range *selection.Directives() {
				if directive.arguments != nil {
					// For each argument, if it is a variable, add it to the variables set.
					// If it is a list, check each item in the list.
					// If it is an object, check each value in the object.
					for _, arg := range *directive.arguments {
						if object, isObject := arg.value.(*objectValue); isObject {
							inputValueDef := getDirectiveArgumentDefinition(schema, directive, arg)

							usages = append(usages, objectValueUsage{
								locationType: inputValueDef._type,
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
			if field.arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.arguments {
					if object, isObject := arg.value.(*objectValue); isObject {
						inputValueDef := getFieldArgumentDefinition(
							schema,
							rootSelectionSet,
							field,
							arg,
							fragmentsPool,
						)

						usages = append(usages, objectValueUsage{
							locationType: inputValueDef._type,
							object:       *object,
						})
					}
				}
			}

			// If the field has a selection set, extract its variables too.
			if field.selectionSet != nil {
				selectionSetObjectUsages := collectObjectValueUsages(
					schema,
					rootSelectionSet,
					*field.selectionSet,
					fragmentsPool,
				)

				usages = append(usages, selectionSetObjectUsages...)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			inlineFragObjectUsages := collectObjectValueUsages(
				schema,
				rootSelectionSet,
				inlineFrag.selectionSet,
				fragmentsPool,
			)

			usages = append(usages, inlineFragObjectUsages...)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			inlineFragObjectUsages := collectObjectValueUsages(
				schema,
				rootSelectionSet,
				fragmentsPool[fragSpread.fragmentName.value].selectionSet,
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
		if selection.Directives() != nil {
			for _, directive := range *selection.Directives() {
				if directive.arguments != nil {
					// For each argument, if it is a variable, add it to the variables set.
					// If it is a list, check each item in the list.
					// If it is an object, check each value in the object.
					for _, arg := range *directive.arguments {
						if _var, isVariable := arg.value.(*variable); isVariable {
							usages[&arg] = &variableUsage{
								field:       nil,
								directive:   directive,
								argument:    arg,
								listValue:   nil,
								objectField: nil,
								variable:    *_var,
							}
						} else if list, isList := arg.value.(*listValue); isList {
							for _, item := range list.values {
								if _var, isVariable := item.(*variable); isVariable {
									usages[&list] = &variableUsage{
										field:       nil,
										directive:   directive,
										argument:    arg,
										listValue:   list,
										objectField: nil,
										variable:    *_var,
									}
								}
							}
						} else if object, isObject := arg.value.(*objectValue); isObject {
							for _, field := range object.values {
								if _var, isVariable := field._value.(*variable); isVariable {
									usages[&object] = &variableUsage{
										field:       nil,
										directive:   directive,
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
			if field.arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.arguments {
					if _var, isVariable := arg.value.(*variable); isVariable {
						usages[&arg] = &variableUsage{
							field:       field,
							directive:   nil,
							argument:    arg,
							listValue:   nil,
							objectField: nil,
							variable:    *_var,
						}
					} else if list, isList := arg.value.(*listValue); isList {
						for _, item := range list.values {
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
					} else if object, isObject := arg.value.(*objectValue); isObject {
						for _, objectField := range object.values {
							if _var, isVariable := objectField._value.(*variable); isVariable {
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
			if field.selectionSet != nil {
				collectVariableUsages(*field.selectionSet, usages, fragmentsPool)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			collectVariableUsages(*inlineFrag.SelectionSet(), usages, fragmentsPool)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			collectVariableUsages(fragmentsPool[fragSpread.fragmentName.value].selectionSet,
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
			targetsSet[fragmentSpread.fragmentName.value] = struct{}{}
		} else if selection.SelectionSet() != nil {
			extractUsedFragmentsNames(*selection.SelectionSet(), targetsSet)
		}
	}
}

func extractUsedVariablesNames(selectionSet selectionSet,
	variablesSet map[string]struct{},
	fragmentsPool map[string]*fragmentDefinition) {
	// Loop over the selection in the selection set.
	for _, selection := range selectionSet {
		// If the selection has directives, check their arguments too.
		if selection.Directives() != nil {
			for _, directive := range *selection.Directives() {
				if directive.arguments != nil {
					// For each argument, if it is a variable, add it to the variables set.
					// If it is a list, check each item in the list.
					// If it is an object, check each value in the object.
					for _, arg := range *directive.arguments {
						if _var, isVariable := arg.value.(*variable); isVariable {
							variablesSet[_var.name.value] = struct{}{}
						} else if list, isList := arg.value.(*listValue); isList {
							for _, item := range list.values {
								if _var, isVariable := item.(*variable); isVariable {
									variablesSet[_var.name.value] = struct{}{}
								}
							}
						} else if object, isObject := arg.value.(*objectValue); isObject {
							for _, field := range object.values {
								if _var, isVariable := field._value.(*variable); isVariable {
									variablesSet[_var.name.value] = struct{}{}
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
			if field.arguments != nil {
				// For each argument, if it is a variable, add it to the variables set.
				// If it is a list, check each item in the list.
				// If it is an object, check each value in the object.
				for _, arg := range *field.arguments {
					if _var, isVariable := arg.value.(*variable); isVariable {
						variablesSet[_var.name.value] = struct{}{}
					} else if list, isList := arg.value.(*listValue); isList {
						for _, item := range list.values {
							if _var, isVariable := item.(*variable); isVariable {
								variablesSet[_var.name.value] = struct{}{}
							}
						}
					} else if object, isObject := arg.value.(*objectValue); isObject {
						for _, field := range object.values {
							if _var, isVariable := field._value.(*variable); isVariable {
								variablesSet[_var.name.value] = struct{}{}
							}
						}
					}
				}
			}

			// If the field has a selection set, extract its variables too.
			if field.selectionSet != nil {
				extractUsedVariablesNames(*field.selectionSet, variablesSet, fragmentsPool)
			}
		} else if inlineFrag, isInlineFrag := selection.(*inlineFragment); isInlineFrag {
			extractUsedVariablesNames(*inlineFrag.SelectionSet(), variablesSet, fragmentsPool)
		} else if fragSpread, isFragSpread := selection.(*fragmentSpread); isFragSpread {
			extractUsedVariablesNames(fragmentsPool[fragSpread.fragmentName.value].selectionSet,
				variablesSet,
				fragmentsPool)
		}
	}
}

func getFragmentsPool(doc *document) map[string]*fragmentDefinition {
	fragmentsPool := make(map[string]*fragmentDefinition)

	// Populate the fragment dictionary with all the fragments in the document.
	// the key of the dictionary is the name of the fragment definition for easy
	// access.
	for _, def := range doc.definitions {
		if fragDef, ok := def.(*fragmentDefinition); ok {
			fragmentsPool[fragDef.fragmentName.value] = fragDef
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
	extractUsedFragmentsNames(fragDef.SelectionSet(), spreads)

	// For each spread, make sure that it is not already exists in visited.
	// If it is, panic.
	for spread := range spreads {
		if _, ok := visited[spread]; ok {
			panic(errors.New("The graph of fragment spreads must not" +
				" form any cycles including spreading itself: \"" + spread + "\"" +
				" has been visited more than once in a single execution path"))
		}

		// Add the spread to the visited set.
		visited[spread] = struct{}{}

		// Call detectFragmentCycles with the target of spread.
		detectFragmentCycles(*fragmentsPool[spread], visited, fragmentsPool)
	}
}

func isInputType(schema *document, variableType _type) bool {
	basicScalars := map[string]struct{}{
		"Boolean": struct{}{},
		"String":  struct{}{},
		"Int":     struct{}{},
		"Float":   struct{}{},
	}

	if _, isBasicType := basicScalars[variableType.TypeName()]; isBasicType {
		return true
	}

	if nonNullType, isNonNullType := variableType.(*nonNullType); isNonNullType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isInputType(schema, nonNullType.ofType)
	}

	if listType, isListType := variableType.(*listType); isListType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isInputType(schema, listType.OfType)
	}

	for _, def := range schema.definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.Name().value == variableType.TypeName() {
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

func isOutputType(schema *document, variableType _type) bool {
	basicScalars := map[string]struct{}{
		"Boolean": struct{}{},
		"String":  struct{}{},
		"Int":     struct{}{},
		"Float":   struct{}{},
	}

	if _, isBasicType := basicScalars[variableType.TypeName()]; isBasicType {
		return true
	}

	if nonNullType, isNonNullType := variableType.(*nonNullType); isNonNullType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isOutputType(schema, nonNullType.ofType)
	}

	if listType, isListType := variableType.(*listType); isListType {
		// Let unwrappedType be the unwrapped type of type.
		// Return IsInputType(unwrappedType)
		return isOutputType(schema, listType.OfType)
	}

	for _, def := range schema.definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.Name().value == variableType.TypeName() {
				switch typeDef.(type) {
				case *scalarTypeDefinition,
					*objectTypeDefinition,
					*interfaceTypeDefinition,
					*unionTypeDefinition,
					*enumTypeDefinition:
					return true
				default:
					return false
				}
			}
		}
	}

	return false
}

func checkIfDirectivesInSelectionSetAreDefined(schema *document, selectionSet selectionSet) {
	// For each selection, look for directives.
	for _, selection := range selectionSet {
		// And the field contains directives
		if selection.Directives() != nil {
			// For each directive, decide whether it is defined or not.
			for _, directive := range *selection.Directives() {
				// If the directive is not defined (which means the directive is not a built in
				// directive and we could not find a proper directive definition in
				// the schema), return false.
				isDirectiveDefined(schema, directive)
			}
		}

		// If the field contains a selection set, check it's directives too.
		if selection.SelectionSet() != nil {
			checkIfDirectivesInSelectionSetAreDefined(schema, *selection.SelectionSet())
		}
	}
}

func isDirectiveDefined(schema *document, directive *directive) {
	// A map that contains all graphql's built in directives for quick access.
	builtInDirectiveNames := map[string]struct{}{
		"skip":       struct{}{},
		"include":    struct{}{},
		"deprecated": struct{}{},
	}

	// Check if the directive is a build in directive. If it is, turn on the
	// flag. Else, Search the directive in the schema.
	if _, isBuiltInDirective := builtInDirectiveNames[directive.name.value]; isBuiltInDirective {
		return
	} else {
		// For each definition in the schema, if it is a directive definition,
		// compare its name to the current directive from the document.
		for _, def := range schema.definitions {
			if directiveDef, isDirectiveDef := def.(*directiveDefinition); isDirectiveDef {
				// If the name are equal, turn on the flag.
				if directiveDef.name.value == directive.name.value {
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
			if field.arguments != nil {
				for _, arg := range *field.arguments {
					if object, isObject := arg.value.(*objectValue); isObject {
						inputObjects = append(inputObjects, object)
					}
				}
			}
		}

		if selection.Directives() != nil {
			for _, directive := range *selection.Directives() {
				if directive.arguments != nil {
					for _, arg := range *directive.arguments {
						if object, isObject := arg.value.(*objectValue); isObject {
							inputObjects = append(inputObjects, object)
						}
					}
				}
			}
		}

		if selection.SelectionSet() != nil {
			inputObjects = append(inputObjects, collectInputObjects(*selection.SelectionSet())...)
		}
	}

	return inputObjects
}

func checkDirectiveLocation(schema *document, usage *directiveUsage) bool {
	// If the directive is one of the built in directives, and it has been used
	// on a field, fragment spread or an inline fragment return true.
	if usage.directive.name.value == "skip" ||
		usage.directive.name.value == "include" {
		if usage.location == edlField || usage.location == edlFragmentSpread ||
			usage.location == edlInlineFragment {
			return true
		}
	}

	//For each definition in the schema:
	for _, def := range schema.definitions {
		// If the definition is a directive definition:
		if dirDef, isDirDef := def.(*directiveDefinition); isDirDef {
			// If the name of the directive definition equals to the name of the
			// given directive:
			if usage.directive.name.value == dirDef.name.value {
				// For each location in the directive definition, compare it with the
				// given location. If they are equal return true.
				for _, location := range dirDef.directiveLocations {
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
		if selection.Directives() != nil {
			// For each directive, set the location to be the kind of the selection.
			for _, dir := range *selection.Directives() {
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
				usages = append(usages, &directiveUsage{dir, location})
			}
		}

		// If the selection has a selection set, get a list of all of its directive usages.
		if selection.SelectionSet() != nil {
			usages = append(usages, extractDirectivesWithLocationsFromSelectionSet(*selection.SelectionSet())...)
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
		if _, isDirNameAlreadyExists := directivesSet[dir.name.value]; isDirNameAlreadyExists {
			return false
		}

		// Insert the directive name to the unique set.
		directivesSet[dir.name.value] = struct{}{}
	}

	return true
}

func checkDirectivesUniquenessInSelectionSet(selectionSet selectionSet) bool {
	// For each selection in the selection set:
	for _, selection := range selectionSet {
		// If the selection has directives, check their uniqueness.
		if selection.Directives() != nil {
			if !checkDirectivesUniqueness(*selection.Directives()) {
				return false
			}
		}

		if selection.SelectionSet() != nil {
			if !checkDirectivesUniquenessInSelectionSet(*selection.SelectionSet()) {
				return false
			}
		}
	}

	return true
}

func getRootQueryTypeDefinition(schema *document) *objectTypeDefinition {
	rootQueryTypeName := "Query"

	for _, def := range schema.definitions {
		if schemaDef, isSchemaDef := def.(*schemaDefinition); isSchemaDef {
			for _, rootOperationType := range schemaDef.rootOperationTypeDefinitions {
				if rootOperationType.operationType == operationQuery {
					rootQueryTypeName = rootOperationType.namedType.value
				}
			}
		}
	}

	for _, def := range schema.definitions {
		if objectTypeDef, isObjectTypeDef := def.(*objectTypeDefinition); isObjectTypeDef {
			if objectTypeDef.name.value == rootQueryTypeName {
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
	schema *document,
	fragmentsPool map[string]*fragmentDefinition,
) _type {
	var typeToReturn _type
	for _, selection := range current {
		switch s := selection.(type) {
		case *field:
			if s.selectionSet != nil {
				var fieldType _type
				switch t := parentType.(type) {
				case *objectTypeDefinition:
					if t.fieldsDefinition != nil {
						for _, fieldDef := range *t.fieldsDefinition {
							if fieldDef.name.value == s.name.value {
								fieldType = fieldDef._type
								break
							}
						}
					}
				case *interfaceTypeDefinition:
					if t.fieldsDefinition != nil {
						for _, fieldDef := range *t.fieldsDefinition {
							if fieldDef.name.value == s.name.value {
								fieldType = fieldDef._type
								break
							}
						}
					}
				case *unionTypeDefinition:
					if t.unionMemberTypes != nil {
						for _, unionMember := range *t.unionMemberTypes {
							fieldType = getSelectionSetType(
								getTypeDefinitionByType(schema, unionMember),
								target,
								*s.selectionSet,
								schema,
								fragmentsPool)

							if fieldType != nil {
								break
							}
						}
					}
				}

				if s.selectionSet == target {
					return fieldType
				} else {
					fieldSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, fieldType),
						target,
						*s.selectionSet,
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
			if s.selectionSet != nil {
				if &s.selectionSet == target {
					return &s.typeCondition.namedType
				} else {
					inlineFragSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, &s.typeCondition.namedType),
						target,
						s.selectionSet,
						schema,
						fragmentsPool,
					)

					if inlineFragSelectionSetType != nil {
						typeToReturn = inlineFragSelectionSetType
					}
				}
			}
		case *fragmentSpread:
			frag := fragmentsPool[s.fragmentName.value]

			if frag.selectionSet != nil {
				if &frag.selectionSet == target {
					return &frag.typeCondition.namedType
				} else {
					fragSelectionSetType := getSelectionSetType(
						getTypeDefinitionByType(schema, &frag.typeCondition.namedType),
						target,
						frag.selectionSet,
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
	schema *document,
	fragmentsPool map[string]*fragmentDefinition,
) *fieldDefinition {
	var tachlessFieldDefinition *fieldDefinition

	for _, selection := range selectionSet {
		switch s := selection.(type) {
		case *field:
			switch t := parentType.(type) {
			case *objectTypeDefinition:
				if t.fieldsDefinition != nil {
					for _, fieldDef := range *t.fieldsDefinition {
						if fieldDef.name.value == s.name.value {
							tachlessFieldDefinition = fieldDef
							break
						}
					}
				}
			case *interfaceTypeDefinition:
				if t.fieldsDefinition != nil {
					for _, fieldDef := range *t.fieldsDefinition {
						if fieldDef.name.value == s.name.value {
							tachlessFieldDefinition = fieldDef
							break
						}
					}
				}
			case *unionTypeDefinition:
				if t.unionMemberTypes != nil {
					for _, unionMember := range *t.unionMemberTypes {
						fieldDef := getFieldDefinitionByFieldSelection(
							getTypeDefinitionByType(schema, unionMember),
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
				panic(errors.New("could not find a field definition named " + s.name.value))
			}

			if selection == targetSelection {
				return tachlessFieldDefinition
			} else if s.SelectionSet() != nil {
				return getFieldDefinitionByFieldSelection(
					getTypeDefinitionByType(schema, tachlessFieldDefinition._type),
					targetSelection,
					*s.SelectionSet(),
					schema,
					fragmentsPool,
				)
			} else {
				return nil
			}
		case *inlineFragment:
			return getFieldDefinitionByFieldSelection(
				getTypeDefinitionByType(schema, &s.typeCondition.namedType),
				targetSelection,
				*s.SelectionSet(),
				schema,
				fragmentsPool,
			)
		case *fragmentSpread:
			return getFieldDefinitionByFieldSelection(
				getTypeDefinitionByType(
					schema,
					&fragmentsPool[s.fragmentName.value].typeCondition.namedType,
				),
				targetSelection,
				*s.SelectionSet(),
				schema, fragmentsPool,
			)
		}
	}

	panic(errors.New("empty selection cannot query for selection type"))
}

func getTypeDefinitionByType(schema *document, t _type) typeDefinition {
	for _, def := range schema.definitions {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.Name().value == t.TypeName() {
				return typeDef
			}
		}
	}

	panic(errors.New("could not find a type definition named: " + t.TypeName()))
}

// http://spec.graphql.org/draft/#GetPossibleTypes()
func getPossibleTypes(schema *document, typeDef typeDefinition) map[string]struct{} {
	typesSet := make(map[string]struct{})

	switch v := typeDef.(type) {
	case *objectTypeDefinition:
		typesSet[v.name.value] = struct{}{}
	case *interfaceTypeDefinition:
		for _, def := range schema.definitions {
			if objectTypeDef, isObjectTypeDef := def.(*objectTypeDefinition); isObjectTypeDef {
				if objectTypeDef.implementsInterfaces != nil {
					for _, iface := range *objectTypeDef.implementsInterfaces {
						if iface.TypeName() == v.name.value {
							typesSet[objectTypeDef.name.value] = struct{}{}
						}
					}
				}
			}
		}
	case *unionTypeDefinition:
		if v.unionMemberTypes != nil {
			for _, unionMember := range *v.unionMemberTypes {
				typesSet[unionMember.TypeName()] = struct{}{}
			}
		}
	default:
		panic(errors.New("Cannot get possible types of a type which is not " +
			"an object, interface or union type definition"))
	}

	return typesSet
}

func checkSpreadsPossibilityInSelectionSet(
	schema *document,
	rootSelectionSet selectionSet,
	selectionSet selectionSet,
	parentType typeDefinition,
	fragmentsPool map[string]*fragmentDefinition,
) {
	// For each selection in the selection set:
	for _, selection := range selectionSet {
		var fragmentType _type

		// Check the selection's type:
		switch s := selection.(type) {
		// If the selection is a field:
		case *field:
			// If the field contains selection set of its own, check its
			// spreads' possibility.
			if s.selectionSet != nil {
				fieldDef := getFieldDefinitionByFieldSelection(
					parentType,
					s,
					rootSelectionSet,
					schema,
					fragmentsPool,
				)

				checkSpreadsPossibilityInSelectionSet(
					schema,
					rootSelectionSet,
					*s.selectionSet,
					getTypeDefinitionByType(schema, fieldDef._type),
					fragmentsPool,
				)
			}
		// If the selection is an inline fragment:
		case *inlineFragment:
			// Save the fragment's type for later use.
			fragmentType = &s.typeCondition.namedType

			// Check the fragment's selection set spreads' possibility.
			checkSpreadsPossibilityInSelectionSet(
				schema,
				rootSelectionSet,
				s.selectionSet,
				getTypeDefinitionByType(schema, fragmentType),
				fragmentsPool,
			)
		// If the selection is a fragment spread:
		case *fragmentSpread:
			// Get the fragment definition that the spread points to.
			fragment := fragmentsPool[s.fragmentName.value]

			// Save the fragment's type for later use.
			fragmentType = &fragment.typeCondition.namedType

			// Check the fragment's selection set spreads' possibility.
			checkSpreadsPossibilityInSelectionSet(
				schema,
				rootSelectionSet,
				fragment.selectionSet,
				getTypeDefinitionByType(schema, &fragment.typeCondition.namedType),
				fragmentsPool,
			)
		}

		// If fragmentType is not nil, it means that the current selection is either
		// an inline fragment or a fragment spread, so we need to check the spread's
		// possibility.
		if fragmentType != nil {
			// Get the possible types of the parent type.
			parentPossibleTypes := getPossibleTypes(schema, parentType)

			// Get the possible types of the fragment type.
			fragmentPossibleTypes := getPossibleTypes(
				schema,
				getTypeDefinitionByType(schema, fragmentType),
			)

			// Create a dictionary that will hold the intersection of the two types
			// dictionaries.
			intersectingTypes := make(map[string]struct{})

			// For each type in the fragment possible types:
			for t := range fragmentPossibleTypes {
				// If it also exists in the parent possible types, insert it to the
				// intersecting types dictionary.
				if _, ok := parentPossibleTypes[t]; ok {
					intersectingTypes[t] = struct{}{}
				}
			}

			// For each type in the parent possible types:
			for t := range parentPossibleTypes {
				// If it also exists in the fragment possible types, insert it to the
				// intersecting types dictionary.
				if _, ok := fragmentPossibleTypes[t]; ok {
					intersectingTypes[t] = struct{}{}
				}
			}

			// If intersectingTypes does not contain any types, the fragment spread
			// is not valid, panic.
			if len(intersectingTypes) < 1 {
				panic(errors.New("A fragment spread is only valid if its type" +
					" condition could ever apply within the parent type."))
			}
		}
	}
}

func isInputFieldDefined(field objectField, fieldsDefinition inputFieldsDefinition) bool {
	for _, fieldDef := range fieldsDefinition {
		if fieldDef.name.value == field.name.value {
			return true
		}
	}

	return false
}

func isIntCoercible(v value) bool {
	if intVal, isInt := v.(*intValue); isInt {
		if intVal._value <= 2^32 || intVal._value >= -(2^32) {
			return true
		}
	}

	return false
}

func isFloatCoercible(v value) bool {
	if _, isFloat := v.(*floatValue); isFloat {
		return true
	} else if _, isInt := v.(*intValue); isInt {
		return true
	}

	return false
}

func isStringCoercible(v value) bool {
	_, isString := v.(*stringValue)
	return isString
}

func isBooleanCoercible(v value) bool {
	_, isBool := v.(*booleanValue)
	return isBool
}

func isIdCoercible(v value) bool {
	if _, isString := v.(*stringValue); isString {
		return true
	} else if _, isInt := v.(*intValue); isInt {
		return true
	}

	return false
}

func isEnumCoercible(v value) bool {
	_, isEnum := v.(*enumValue)
	return isEnum
}

func isInputObjectCoercible(v value) bool {
	// TODO: Implement
	return false
}
