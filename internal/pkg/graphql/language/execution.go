package language

import "fmt"

func collectFields(doc *document, objectType *rootOperationTypeDefinition, selectionSet selectionSet, variableValues []value, visitedFragments ...*name) map[string][]selection {
	groupedFields := make(map[string][]selection)

	for _, _selection := range selectionSet {
		// Check @skip directive
		if dirExists, index := execDirectiveExists(*_selection.Directives(), "skip"); dirExists {
			skipDirective := (*_selection.Directives())[index]

			if skipDirective.Arguments() != nil {
				if argExists, index := execArgumentExists(*skipDirective.Arguments(), "if"); argExists {
					ifArg := (*skipDirective.Arguments())[index]

					// This case should handle variable values as well
					if val, ok := ifArg.Value().Value().(bool); ok {
						if val {
							continue
						}
					}
				}
			}
		}

		// Check @include directive
		if dirExists, index := execDirectiveExists(*_selection.Directives(), "include"); dirExists {
			includeDirective := (*_selection.Directives())[index]

			if includeDirective.arguments != nil {
				if argExists, index := execArgumentExists(*includeDirective.Arguments(), "if"); argExists {
					ifArg := (*includeDirective.Arguments())[index]

					// This case should handle variable values as well
					if val, ok := ifArg.Value().Value().(bool); ok {
						if !val {
							continue
						}
					}
				}
			}
		}

		// If `selection` is a field
		if field, isField := _selection.(*field); isField {
			var responseKey string

			if field.alias != nil {
				responseKey = field.Alias().Value()
			} else {
				name := field.Name()
				responseKey = name.Value()
			}

			_, exists := groupedFields[responseKey]

			if !exists {
				groupedFields[responseKey] = make([]selection, 0)
			}

			groupedFields[responseKey] = append(groupedFields[responseKey], _selection)

			continue
		}

		// If `selection` is a fragment spread
		if fs, isFs := _selection.(*fragmentSpread); isFs {
			fragmentSpreadName := fs.FragmentName()

			if visitedFragmentsContainFragmentName(visitedFragments, fragmentSpreadName.Value()) {
				continue
			}

			visitedFragments = append(visitedFragments, fs.FragmentName())

			fragment := getFragmentByName(doc, fragmentSpreadName.Value())

			if fragment == nil {
				continue
			}

			fragmentType := fragment.TypeCondition()

			if !DoesFragmentTypeApply(doc, objectType, fragmentType) {
				continue
			}

		}
	}

	return groupedFields
}

func DoesFragmentTypeApply(doc *document, objectType *rootOperationTypeDefinition, fragmentType *typeCondition) bool {
	typeDefinition := getTypeDefinitionByName(doc, fragmentType.NamedType().TypeName())

	if objectType, isObjectType := typeDefinition.(*objectTypeDefinition); isObjectType {
		return objectType.Name().Value() == fragmentType.NamedType().TypeName()
	}

	// if interfaceType, isInterfaceType := typeDefinition.(*interfaceTypeDefinition); isInterfaceType {
	// 	objectTypeDefinition := getTypeDefinitionByName(doc, objectType.NamedType().TypeName())

	// 	// for _, implementedInterface := objectTypeDefinition.ImplementsInterfaces() {

	// 	// }
	// }

	return false
}

// Helper Functions
func getTypeDefinitionByName(doc *document, namedType string) typeDefinition {
	for _, def := range doc.Definitions() {
		if typeDef, isTypeDef := def.(typeDefinition); isTypeDef {
			if typeDef.Name().Value() == namedType {
				return typeDef
			}
		}
	}

	panic(fmt.Sprintf("no type found by name %s", namedType))
}

func getFragmentByName(doc *document, fragmentName string) *fragmentDefinition {
	for _, def := range doc.Definitions() {
		if fragmentDef, isFragmentDef := def.(*fragmentDefinition); isFragmentDef {
			if fragmentDef.Name().Value() == fragmentName {
				return fragmentDef
			}
		}
	}

	return nil
}

func execDirectiveExists(dirs directives, dirName string) (bool, int) {
	for i, dir := range dirs {
		name := dir.Name()
		if name.Value() == dirName {
			return true, i
		}
	}

	return false, -1
}

func execArgumentExists(args arguments, argName string) (bool, int) {
	for i, arg := range args {
		name := arg.Name()
		if name.Value() == argName {
			return true, i
		}
	}

	return false, -1
}

func visitedFragmentsContainFragmentName(visitedFragments []*name, fragName string) bool {
	for _, fs := range visitedFragments {
		if fs.Value() == fragName {
			return true
		}
	}

	return false
}
