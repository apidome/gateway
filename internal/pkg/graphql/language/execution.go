package language

func collectFields(objectType interface{}, selectionSet *selectionSet, variableValues interface{}, visitedFragments *[]fragmentSpread) map[string][]interface{} {
	if visitedFragments == nil {
		vf := make([]fragmentSpread, 0)

		visitedFragments = &vf
	}
	groupedFields := make(map[string][]interface{})

	return groupedFields
}
