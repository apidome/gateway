package language

func collectFields(objectType interface{}, selectionSet *selectionSet, variableValues map[string]interface{}, visitedFragments *[]fragmentSpread) map[string][]interface{} {
	if visitedFragments == nil {
		vf := make([]fragmentSpread, 0)

		visitedFragments = &vf
	}

	groupedFields := make(map[string][]interface{})

	for _, selection := range *selectionSet {
		if dirExists, index := execDirectiveExists(*selection.GetDirectives(), "skip"); dirExists {
			skipDirective := (*selection.GetDirectives())[index]

			if skipDirective.Arguments != nil {
				if argExists, index := execArgumentExists(*skipDirective.Arguments, "if"); argExists {
					ifArg := (*skipDirective.Arguments)[index]

					if val, ok := ifArg.Value.GetValue().(bool); ok {
						if val {
							continue
						}
					}
				}
			}
		}
	}

	return groupedFields
}

// Helper Functions
func execDirectiveExists(dirs directives, dirName string) (bool, int) {
	for i, dir := range dirs {
		if dir.Name.Value == dirName {
			return true, i
		}
	}

	return false, -1
}

func execArgumentExists(args arguments, argName string) (bool, int) {
	for i, arg := range args {
		if arg.Name.Value == argName {
			return true, i
		}
	}

	return false, -1
}
