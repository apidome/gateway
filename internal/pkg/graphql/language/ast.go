package language

import "encoding/json"

const (
	// Executable Directive Locations
	edlQuery              executableDirectiveLocation = "QUERY"
	edlMutation                                       = "MUTATION"
	edlSubscription                                   = "SUBSCRIPTION"
	edlField                                          = "FIELD"
	edlFragmentDefinition                             = "FRAGMENT_DEFINITION"
	edlFragmentSpread                                 = "FRAGMENT_SPREAD"
	edlInlineFragment                                 = "INLINE_FRAGMENT"
	edlVariableDefinition                             = "VARIABLE_DEFINITION"

	// Type System Directive Locations
	tsdlSchema               typeSystemDirectiveLocation = "SCHEMA"
	tsdlScalar                                           = "SCALAR"
	tsdlObject                                           = "OBJECT"
	tsdlFieldDefinition                                  = "FIELD_DEFINITION"
	tsdlArgumentDefinition                               = "ARGUMENT_DEFINITION"
	tsdlInterface                                        = "INTERFACE"
	tsdlUnion                                            = "UNION"
	tsdlEnum                                             = "ENUM"
	tsdlEnumValue                                        = "ENUM_VALUE"
	tsdlInputObject                                      = "INPUT_OBJECT"
	tsdlInputFieldDefinition                             = "INPUT_FIELD_DEFINITION"
)

const (
	operationQuery        operationType = "query"
	operationMutation     operationType = "mutation"
	operationSubscription operationType = "subscription"
)

var executableDirectiveLocations []executableDirectiveLocation = []executableDirectiveLocation{
	edlQuery,
	edlMutation,
	edlSubscription,
	edlField,
	edlFragmentDefinition,
	edlFragmentSpread,
	edlInlineFragment,
	edlVariableDefinition,
}

var typeSystemDirectiveLocations []typeSystemDirectiveLocation = []typeSystemDirectiveLocation{
	tsdlSchema,
	tsdlScalar,
	tsdlObject,
	tsdlFieldDefinition,
	tsdlArgumentDefinition,
	tsdlInterface,
	tsdlUnion,
	tsdlEnum,
	tsdlEnumValue,
	tsdlInputObject,
	tsdlInputFieldDefinition,
}

type directiveLocation string

type executableDirectiveLocation directiveLocation

type typeSystemDirectiveLocation directiveLocation

type description stringValue

type location struct {
	start  int
	end    int
	source string
}

func (l *location) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Start  int    `json:"start"`
		End    int    `json:"end"`
		Source string `json:"source"`
	}{
		Start:  l.start,
		End:    l.end,
		Source: l.source,
	})
}

func (loc *location) Start() int {
	return loc.start
}

func (loc *location) End() int {
	return loc.end
}

func (loc *location) Source() string {
	return loc.source
}

type document struct {
	definitions []definition
}

func (d *document) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Definitions []definition `json:"definitions"`
	}{
		Definitions: d.definitions,
	})
}

func (doc *document) Definitions() []definition {
	return doc.definitions
}

type arguments []*argument

type argument struct {
	name  name
	value value
	locator
}

func (a *argument) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    name    `json:"name"`
		Value   value   `json:"value"`
		Locator locator `json:"location"`
	}{
		Name:    a.name,
		Value:   a.value,
		Locator: a.locator,
	})
}

func (arg *argument) Name() name {
	return arg.name
}

func (arg *argument) Value() value {
	return arg.value
}

type directiveLocations []directiveLocation

type directives []*directive

type directive struct {
	name      name
	arguments *arguments
	locator
}

func (d *directive) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      name       `json:"name"`
		Arguments *arguments `json:"arguments,omitempty"`
		Locator   locator    `json:"location"`
	}{
		Name:      d.name,
		Arguments: d.arguments,
		Locator:   d.locator,
	})
}

func (dir *directive) Name() name {
	return dir.name
}

func (dir *directive) Arguments() *arguments {
	return dir.arguments
}

type operationType string

type locatorInterface interface {
	Location() *location
}

type locator struct {
	loc location
}

func (l *locator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Location location `json:"location"`
	}{
		Location: l.loc,
	})
}

func (l *locator) Location() *location {
	return &l.loc
}

type definition interface {
	locatorInterface
	definition() definition
}

type definitions []definition

var _ definition = (executableDefinition)(nil)
var _ definition = (typeSystemDefinition)(nil)
var _ definition = (typeSystemExtension)(nil)

/*
	Executable Definitions
*/
type executableDefinition interface {
	definition
	executableDefinition() executableDefinition
	Name() *name
	Directives() *directives
	SelectionSet() selectionSet
}

var _ executableDefinition = (*operationDefinition)(nil)
var _ executableDefinition = (*fragmentDefinition)(nil)

type operationDefinition struct {
	operationType       operationType
	name                *name
	variableDefinitions *variableDefinitions
	selectionSet        selectionSet
	directives          *directives
	locator
}

func (o *operationDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OperationType       operationType        `json:"operationType"`
		Name                *name                `json:"name,omitempty"`
		VariableDefinitions *variableDefinitions `json:"variableDefinitions,omitempty"`
		SelectionSet        selectionSet         `json:"selectionSet"`
		Directives          *directives          `json:"directives,omitempty"`
		Locator             locator              `json:"location"`
	}{
		OperationType:       o.operationType,
		Name:                o.name,
		VariableDefinitions: o.variableDefinitions,
		SelectionSet:        o.selectionSet,
		Directives:          o.directives,
		Locator:             o.locator,
	})
}

func (o *operationDefinition) definition() definition {
	return o
}

func (o *operationDefinition) executableDefinition() executableDefinition {
	return o
}

func (o *operationDefinition) Name() *name {
	return o.name
}

func (o *operationDefinition) Directives() *directives {
	return o.directives
}

func (o *operationDefinition) SelectionSet() selectionSet {
	return o.selectionSet
}

func (o *operationDefinition) OperationType() operationType {
	return o.operationType
}

func (o *operationDefinition) VariableDefinitions() *variableDefinitions {
	return o.variableDefinitions
}

type fragmentDefinition struct {
	fragmentName  name
	typeCondition typeCondition
	selectionSet  selectionSet
	directives    *directives
	locator
}

func (f *fragmentDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FragmentName  name          `json:"name"`
		TypeCondition typeCondition `json:"typeCondition"`
		SelectionSet  selectionSet  `json:"selectionSet"`
		Directives    *directives   `json:"directives,omitempty"`
		Locator       locator       `json:"location"`
	}{
		FragmentName:  f.fragmentName,
		TypeCondition: f.typeCondition,
		SelectionSet:  f.selectionSet,
		Directives:    f.directives,
		Locator:       f.locator,
	})
}

func (f *fragmentDefinition) definition() definition {
	return f
}

func (f *fragmentDefinition) executableDefinition() executableDefinition {
	return f
}

func (f *fragmentDefinition) Name() *name {
	return &f.fragmentName
}

func (f *fragmentDefinition) Directives() *directives {
	return f.directives
}

func (f *fragmentDefinition) SelectionSet() selectionSet {
	return f.selectionSet
}

func (f *fragmentDefinition) TypeCondition() typeCondition {
	return f.typeCondition
}

/*
	Type System Definitions
*/

type typeSystemDefinition interface {
	definition
	typeSystemDefinition() typeSystemDefinition
}

var _ typeSystemDefinition = (typeDefinition)(nil)
var _ typeSystemDefinition = (*schemaDefinition)(nil)
var _ typeSystemDefinition = (*directiveDefinition)(nil)

type schemaDefinition struct {
	description                  *description
	directives                   *directives
	rootOperationTypeDefinitions rootOperationTypeDefinitions
	locator
}

func (s *schemaDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description                  *description                 `json:"description,omitmepty"`
		Directives                   *directives                  `json:"directives,omitempty"`
		RootOperationTypeDefinitions rootOperationTypeDefinitions `json:"rootOperationTypeDefinitions"`
		Locator                      locator                      `json:"location"`
	}{
		Description:                  s.description,
		Directives:                   s.directives,
		RootOperationTypeDefinitions: s.rootOperationTypeDefinitions,
		Locator:                      s.locator,
	})
}

func (s *schemaDefinition) definition() definition {
	return s
}

func (s *schemaDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *schemaDefinition) Directives() *directives {
	return s.directives
}

func (s *schemaDefinition) Description() *description {
	return s.description
}

func (s *schemaDefinition) RootOperationTypeDefinitions() rootOperationTypeDefinitions {
	return s.rootOperationTypeDefinitions
}

type rootOperationTypeDefinition struct {
	operationType operationType
	namedType     namedType
	locator
}

func (r *rootOperationTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OperationType operationType `json:"operationType"`
		NamedType     namedType     `json:"namedType"`
		Locator       locator       `json:"location"`
	}{
		OperationType: r.operationType,
		NamedType:     r.namedType,
		Locator:       r.locator,
	})
}

func (rotd *rootOperationTypeDefinition) OperationType() operationType {
	return rotd.operationType
}

func (rotd *rootOperationTypeDefinition) NamedType() namedType {
	return rotd.namedType
}

type rootOperationTypeDefinitions []*rootOperationTypeDefinition

type directiveDefinition struct {
	description         *description
	name                name
	argumentsDefinition *argumentsDefinition
	directiveLocations  directiveLocations
	locator
}

func (d *directiveDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description         *description         `json:"description"`
		Name                name                 `json:"name"`
		ArgumentsDefinition *argumentsDefinition `json:"argumentsDefinition,omitempty"`
		DirectiveLocations  directiveLocations   `json:"directiveLocations"`
		Locator             locator              `json:"location"`
	}{
		Description:         d.description,
		Name:                d.name,
		ArgumentsDefinition: d.argumentsDefinition,
		DirectiveLocations:  d.directiveLocations,
		Locator:             d.locator,
	})
}

func (d *directiveDefinition) definition() definition {
	return d
}

func (d *directiveDefinition) typeSystemDefinition() typeSystemDefinition {
	return d
}

func (dd *directiveDefinition) Description() *description {
	return dd.description
}

func (dd *directiveDefinition) Name() name {
	return dd.name
}

func (dd *directiveDefinition) ArgumentsDefinition() *argumentsDefinition {
	return dd.argumentsDefinition
}

func (dd *directiveDefinition) DirectiveLocation() directiveLocations {
	return dd.directiveLocations
}

type inputValueDefinition struct {
	description  *description
	name         name
	_type        _type
	defaultValue *defaultValue
	directives   *directives
	locator
}

func (i *inputValueDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description  *description  `json:"description,omitempty"`
		Name         name          `json:"name"`
		Type         _type         `json:"type"`
		DefaultValue *defaultValue `json:"defaultValue,omitempty"`
		Directives   *directives   `json:"directives,omitempty"`
		Locator      locator       `json:"location"`
	}{
		Description:  i.description,
		Name:         i.name,
		Type:         i._type,
		DefaultValue: i.defaultValue,
		Directives:   i.directives,
		Locator:      i.locator,
	})
}

func (ivd *inputValueDefinition) Description() *description {
	return ivd.description
}

func (ivd *inputValueDefinition) Name() name {
	return ivd.name
}

func (ivd *inputValueDefinition) Type() _type {
	return ivd._type
}

func (ivd *inputValueDefinition) DefaultValue() *defaultValue {
	return ivd.defaultValue
}

func (ivd *inputValueDefinition) Directives() *directives {
	return ivd.directives
}

type argumentsDefinition []*inputValueDefinition

type typeDefinition interface {
	typeSystemDefinition
	typeDefinition() typeDefinition
	Description() *description
	Name() name
	Directives() *directives
}

var _ typeDefinition = (*scalarTypeDefinition)(nil)
var _ typeDefinition = (*objectTypeDefinition)(nil)
var _ typeDefinition = (*interfaceTypeDefinition)(nil)
var _ typeDefinition = (*unionTypeDefinition)(nil)
var _ typeDefinition = (*enumTypeDefinition)(nil)
var _ typeDefinition = (*inputObjectTypeDefinition)(nil)

type scalarTypeDefinition struct {
	description *description
	name        name
	directives  *directives
	locator
}

func (s *scalarTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description *description `json:"description,omitempty"`
		Name        name         `json:"name"`
		Directives  *directives  `json:"directives,omitempty"`
		Locator     locator      `json:"location"`
	}{
		Description: s.description,
		Name:        s.name,
		Directives:  s.directives,
		Locator:     s.locator,
	})
}

func (s *scalarTypeDefinition) definition() definition {
	return s
}

func (s *scalarTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *scalarTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (s *scalarTypeDefinition) Description() *description {
	return s.description
}

func (s *scalarTypeDefinition) Name() name {
	return s.name
}

func (s *scalarTypeDefinition) Directives() *directives {
	return s.directives
}

type objectTypeDefinition struct {
	description          *description
	name                 name
	implementsInterfaces *implementsInterfaces
	directives           *directives
	fieldsDefinition     *fieldsDefinition
	locator
}

func (o *objectTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description          *description          `json:"description,omitempty"`
		Name                 name                  `json:"name"`
		ImplementsInterfaces *implementsInterfaces `json:"implementsInterfaces,omitempty"`
		Directives           *directives           `json:"directives,omitempty"`
		FieldsDefinition     *fieldsDefinition     `json:"fieldsDefinition,omitempty"`
		Locator              locator               `json:"location"`
	}{
		Description:          o.description,
		Name:                 o.name,
		ImplementsInterfaces: o.implementsInterfaces,
		Directives:           o.directives,
		FieldsDefinition:     o.fieldsDefinition,
		Locator:              o.locator,
	})
}

func (s *objectTypeDefinition) definition() definition {
	return s
}

func (s *objectTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *objectTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (o *objectTypeDefinition) Description() *description {
	return o.description
}

func (o *objectTypeDefinition) Name() name {
	return o.name
}

func (o *objectTypeDefinition) Directives() *directives {
	return o.directives
}

func (o *objectTypeDefinition) ImplementsInterfaces() *implementsInterfaces {
	return o.implementsInterfaces
}

func (o *objectTypeDefinition) FieldsDefinition() *fieldsDefinition {
	return o.fieldsDefinition
}

type interfaceTypeDefinition struct {
	description      *description
	name             name
	directives       *directives
	fieldsDefinition *fieldsDefinition
	locator
}

func (i *interfaceTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description      *description      `json:"description,omitempty"`
		Name             name              `json:"name"`
		Directives       *directives       `json:"directives,omitempty"`
		FieldsDefinition *fieldsDefinition `json:"fieldsDefinition,omitempty"`
		Locator          locator           `json:"location"`
	}{
		Description:      i.description,
		Name:             i.name,
		Directives:       i.directives,
		FieldsDefinition: i.fieldsDefinition,
		Locator:          i.locator,
	})
}

func (s *interfaceTypeDefinition) definition() definition {
	return s
}

func (s *interfaceTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *interfaceTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (i *interfaceTypeDefinition) Description() *description {
	return i.description
}

func (i *interfaceTypeDefinition) Name() name {
	return i.name
}

func (i *interfaceTypeDefinition) Directives() *directives {
	return i.directives
}

func (i *interfaceTypeDefinition) FieldsDefinition() *fieldsDefinition {
	return i.fieldsDefinition
}

type unionTypeDefinition struct {
	description      *description
	name             name
	directives       *directives
	unionMemberTypes *unionMemberTypes
	locator
}

func (u *unionTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description      *description      `json:"description,omitempty"`
		Name             name              `json:"name"`
		Directives       *directives       `json:"directives,omitempty"`
		UnionMemberTypes *unionMemberTypes `json:"unionMemberTypes,omitempty"`
		Location         locator           `json:"location"`
	}{
		Description:      u.description,
		Name:             u.name,
		Directives:       u.directives,
		UnionMemberTypes: u.unionMemberTypes,
		Location:         u.locator,
	})
}

func (s *unionTypeDefinition) definition() definition {
	return s
}

func (s *unionTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *unionTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (u *unionTypeDefinition) Description() *description {
	return u.description
}

func (u *unionTypeDefinition) Name() name {
	return u.name
}

func (u *unionTypeDefinition) Directives() *directives {
	return u.directives
}

func (u *unionTypeDefinition) UnionMemberTypes() *unionMemberTypes {
	return u.unionMemberTypes
}

type enumTypeDefinition struct {
	description          *description
	name                 name
	directives           *directives
	enumValuesDefinition *enumValuesDefinition
	locator
}

func (e *enumTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description          *description          `json:"description,omitmepty"`
		Name                 name                  `json:"name"`
		Directives           *directives           `json:"directives,omitempty"`
		EnumValuesDefinition *enumValuesDefinition `json:"enumValuesDefinition,omitempty"`
		Locator              locator               `json:"location"`
	}{
		Description:          e.description,
		Name:                 e.name,
		Directives:           e.directives,
		EnumValuesDefinition: e.enumValuesDefinition,
		Locator:              e.locator,
	})
}

func (s *enumTypeDefinition) definition() definition {
	return s
}

func (s *enumTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *enumTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (e *enumTypeDefinition) Description() *description {
	return e.description
}

func (e *enumTypeDefinition) Name() name {
	return e.name
}

func (e *enumTypeDefinition) Directives() *directives {
	return e.directives
}

func (e *enumTypeDefinition) EnumValuesDefinition() *enumValuesDefinition {
	return e.enumValuesDefinition
}

type enumValueDefinition struct {
	description *description
	enumValue   enumValue
	directives  *directives
	locator
}

func (e *enumValueDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description *description `json:"description,omitempty"`
		EnumValue   enumValue    `json:"enumValue"`
		Directives  *directives  `json:"directives,omitempty"`
		Locator     locator      `json:"location"`
	}{
		Description: e.description,
		EnumValue:   e.enumValue,
		Directives:  e.directives,
		Locator:     e.locator,
	})
}

func (e *enumValueDefinition) Description() *description {
	return e.description
}

func (e *enumValueDefinition) EnumValue() enumValue {
	return e.enumValue
}

func (e *enumValueDefinition) Directives() *directives {
	return e.directives
}

type enumValuesDefinition []*enumValueDefinition

type inputObjectTypeDefinition struct {
	description           *description
	name                  name
	directives            *directives
	inputFieldsDefinition *inputFieldsDefinition
	locator
}

func (i *inputObjectTypeDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description           *description           `json:"description,omitempty"`
		Name                  name                   `json:"name"`
		Directives            *directives            `json:"directives,omitempty"`
		InputFieldsDefinition *inputFieldsDefinition `json:"inputFieldsDefinition,omitempty"`
		Locator               locator                `json:"location"`
	}{
		Description:           i.description,
		Name:                  i.name,
		Directives:            i.directives,
		InputFieldsDefinition: i.inputFieldsDefinition,
		Locator:               i.locator,
	})
}

func (s *inputObjectTypeDefinition) definition() definition {
	return s
}

func (s *inputObjectTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *inputObjectTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (i *inputObjectTypeDefinition) Description() *description {
	return i.description
}

func (i *inputObjectTypeDefinition) Name() name {
	return i.name
}

func (i *inputObjectTypeDefinition) Directives() *directives {
	return i.directives
}

func (i *inputObjectTypeDefinition) InputFieldsDefinition() *inputFieldsDefinition {
	return i.inputFieldsDefinition
}

type inputFieldsDefinition []*inputValueDefinition

type variableDefinition struct {
	variable     variable
	_type        _type
	defaultValue *defaultValue
	directives   *directives
	locator
}

func (v *variableDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Variable     variable      `json:"variable"`
		Type         _type         `json:"type"`
		DefaultValue *defaultValue `json:"defaultValue,omitempty"`
		Directives   *directives   `json:"directives,omitempty"`
		Locator      locator       `json:"location"`
	}{
		Variable:     v.variable,
		Type:         v._type,
		DefaultValue: v.defaultValue,
		Directives:   v.directives,
		Locator:      v.locator,
	})
}

func (v *variableDefinition) Variable() variable {
	return v.variable
}

func (v *variableDefinition) Type() _type {
	return v._type
}

func (v *variableDefinition) DefaultValue() *defaultValue {
	return v.defaultValue
}

func (v *variableDefinition) Directives() *directives {
	return v.directives
}

type variableDefinitions []*variableDefinition

type fieldDefinition struct {
	description         *description
	name                name
	argumentsDefinition *argumentsDefinition
	_type               _type
	directives          *directives
	locator
}

func (f *fieldDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description         *description         `json:"description,omitempty"`
		Name                name                 `json:"name"`
		ArgumentsDefinition *argumentsDefinition `json:"argumentsDefinition,omitempty"`
		Type                _type                `json:"type"`
		Directives          *directives          `json:"directives,omitempty"`
		Locator             locator              `json:"location"`
	}{
		Description:         f.description,
		Name:                f.name,
		ArgumentsDefinition: f.argumentsDefinition,
		Type:                f._type,
		Directives:          f.directives,
		Locator:             f.locator,
	})
}

func (f *fieldDefinition) Description() *description {
	return f.description
}

func (f *fieldDefinition) Name() name {
	return f.name
}

func (f *fieldDefinition) ArgumentsDefinition() *argumentsDefinition {
	return f.argumentsDefinition
}

func (f *fieldDefinition) Type() _type {
	return f._type
}

func (f *fieldDefinition) Directives() *directives {
	return f.directives
}

type fieldsDefinition []*fieldDefinition

/*
	Type System Extensions
*/

type typeSystemExtension interface {
	definition
	typeSystemExtension() typeSystemExtension
	Directives() *directives
}

var _ typeSystemExtension = (*schemaExtension)(nil)
var _ typeSystemExtension = (typeExtension)(nil)

type schemaExtension struct {
	directives                   *directives
	rootOperationTypeDefinitions *rootOperationTypeDefinitions
	locator
}

func (s *schemaExtension) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Directives                   *directives                   `json:"directives,omitempty"`
		RootOperationTypeDefinitions *rootOperationTypeDefinitions `json:"rootOperationTypeDefinitions,omitempty"`
		Locator                      locator                       `json:"location"`
	}{
		Directives:                   s.directives,
		RootOperationTypeDefinitions: s.rootOperationTypeDefinitions,
		Locator:                      s.locator,
	})
}

func (s *schemaExtension) definition() definition {
	return s
}

func (s *schemaExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *schemaExtension) Directives() *directives {
	return s.directives
}

func (s *schemaExtension) RootOperationTypeDefinitions() *rootOperationTypeDefinitions {
	return s.rootOperationTypeDefinitions
}

type typeExtension interface {
	typeSystemExtension
	typeExtension() typeExtension
	Name() name
}

var _ typeExtension = (*scalarTypeExtension)(nil)
var _ typeExtension = (*objectTypeExtension)(nil)
var _ typeExtension = (*interfaceTypeExtension)(nil)
var _ typeExtension = (*unionTypeExtension)(nil)
var _ typeExtension = (*enumTypeExtension)(nil)
var _ typeExtension = (*inputObjectTypeExtension)(nil)

type scalarTypeExtension struct {
	name       name
	directives directives
	locator
}

func (s *scalarTypeExtension) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       name       `json:"name"`
		Directives directives `json:"directives"`
		Locator    locator    `json:"location"`
	}{
		Name:       s.name,
		Directives: s.directives,
		Locator:    s.locator,
	})
}

func (s *scalarTypeExtension) definition() definition {
	return s
}

func (s *scalarTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *scalarTypeExtension) typeExtension() typeExtension {
	return s
}

func (s *scalarTypeExtension) Name() name {
	return s.name
}

func (s *scalarTypeExtension) Directives() *directives {
	return &s.directives
}

type objectTypeExtension struct {
	name                 name
	implementsInterfaces *implementsInterfaces
	directives           *directives
	fieldsDefinition     *fieldsDefinition
	locator
}

func (o *objectTypeExtension) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name                 name                  `json:"name"`
		ImplementsInterfaces *implementsInterfaces `json:"implementsInterfaces,omitempty"`
		Directives           *directives           `json:"directives,omitempty"`
		FieldsDefinition     *fieldsDefinition     `json:"fieldsDefinition,omitempty"`
		Locator              locator               `json:"location"`
	}{
		Name:                 o.name,
		ImplementsInterfaces: o.implementsInterfaces,
		Directives:           o.directives,
		FieldsDefinition:     o.fieldsDefinition,
		Locator:              o.locator,
	})
}

func (s *objectTypeExtension) definition() definition {
	return s
}

func (s *objectTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *objectTypeExtension) typeExtension() typeExtension {
	return s
}

func (o *objectTypeExtension) Name() name {
	return o.name
}

func (o *objectTypeExtension) Directives() *directives {
	return o.directives
}

func (o *objectTypeExtension) ImplementsInterfaces() *implementsInterfaces {
	return o.implementsInterfaces
}

func (o *objectTypeExtension) FieldsDefinition() *fieldsDefinition {
	return o.fieldsDefinition
}

type interfaceTypeExtension struct {
	name             name
	directives       *directives
	fieldsDefinition *fieldsDefinition
	locator
}

func (i *interfaceTypeExtension) definition() definition {
	return i
}

func (i *interfaceTypeExtension) typeSystemExtension() typeSystemExtension {
	return i
}

func (i *interfaceTypeExtension) typeExtension() typeExtension {
	return i
}

func (i *interfaceTypeExtension) Name() name {
	return i.name
}

func (i *interfaceTypeExtension) Directives() *directives {
	return i.directives
}

func (i *interfaceTypeExtension) FieldsDefinition() *fieldsDefinition {
	return i.fieldsDefinition
}

type unionTypeExtension struct {
	name             name
	directives       *directives
	unionMemberTypes *unionMemberTypes
	locator
}

func (s *unionTypeExtension) definition() definition {
	return s
}

func (s *unionTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *unionTypeExtension) typeExtension() typeExtension {
	return s
}

func (u *unionTypeExtension) Name() name {
	return u.name
}

func (u *unionTypeExtension) Directives() *directives {
	return u.directives
}

func (u *unionTypeExtension) UnionMemberTypes() *unionMemberTypes {
	return u.unionMemberTypes
}

type enumTypeExtension struct {
	name                 name
	directives           *directives
	enumValuesDefinition *enumValuesDefinition
	locator
}

func (s *enumTypeExtension) definition() definition {
	return s
}

func (s *enumTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *enumTypeExtension) typeExtension() typeExtension {
	return s
}

func (e *enumTypeExtension) Name() name {
	return e.name
}

func (e *enumTypeExtension) Directives() *directives {
	return e.directives
}

func (e *enumTypeExtension) EnumValuesDefinition() *enumValuesDefinition {
	return e.enumValuesDefinition
}

type inputObjectTypeExtension struct {
	name                  name
	directives            *directives
	inputFieldsDefinition *inputFieldsDefinition
	locator
}

func (s *inputObjectTypeExtension) definition() definition {
	return s
}

func (s *inputObjectTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *inputObjectTypeExtension) typeExtension() typeExtension {
	return s
}

func (i *inputObjectTypeExtension) Name() name {
	return i.name
}

func (i *inputObjectTypeExtension) Directives() *directives {
	return i.directives
}

func (i *inputObjectTypeExtension) InputFieldsDefinition() *inputFieldsDefinition {
	return i.inputFieldsDefinition
}

type selection interface {
	locatorInterface
	selection() selection
	SelectionSet() *selectionSet
	Directives() *directives
}

var _ selection = (*field)(nil)
var _ selection = (*fragmentSpread)(nil)
var _ selection = (*inlineFragment)(nil)

type selectionSet []selection

type field struct {
	alias        *alias
	name         name
	arguments    *arguments
	directives   *directives
	selectionSet *selectionSet
	locator
}

func (f *field) selection() selection {
	return f
}

func (f *field) SelectionSet() *selectionSet {
	return f.selectionSet
}

func (f *field) Directives() *directives {
	return f.directives
}

func (f *field) Alias() *alias {
	return f.alias
}

func (f *field) Name() name {
	return f.name
}

func (f *field) Arguments() *arguments {
	return f.arguments
}

type fragmentSpread struct {
	fragmentName name
	directives   *directives
	locator
}

func (f *fragmentSpread) selection() selection {
	return f
}

// TODO Find the fragment definition by its name and return its selection set
func (f *fragmentSpread) SelectionSet() *selectionSet {
	panic("not implemented")
	return &selectionSet{}
}

func (f *fragmentSpread) Directives() *directives {
	return f.directives
}

func (f *fragmentSpread) FragmentName() name {
	return f.fragmentName
}

type inlineFragment struct {
	typeCondition *typeCondition
	directives    *directives
	selectionSet  selectionSet
	locator
}

func (i *inlineFragment) selection() selection {
	return i
}

func (i *inlineFragment) SelectionSet() *selectionSet {
	return &i.selectionSet
}

func (i *inlineFragment) Directives() *directives {
	return i.directives
}

func (i *inlineFragment) TypeCondition() *typeCondition {
	return i.typeCondition
}

type _type interface {
	_type() _type
	TypeName() string
	locatorInterface
}

var _ _type = (*namedType)(nil)
var _ _type = (*listType)(nil)
var _ _type = (*nonNullType)(nil)

type namedType struct {
	name
}

func (n *namedType) _type() _type {
	return n
}

func (n *namedType) TypeName() string {
	return n.value
}

type listType struct {
	kind   string
	OfType _type
	locator
}

func (n *listType) _type() _type {
	return n
}

func (n *listType) TypeName() string {
	return n.OfType.TypeName()
}

func (l *listType) Kind() string {
	return l.kind
}

type nonNullType struct {
	kind   string
	ofType _type
	locator
}

func (n *nonNullType) _type() _type {
	return n
}

func (n *nonNullType) TypeName() string {
	return n.ofType.TypeName()
}

func (n *nonNullType) Kind() string {
	return n.kind
}

type typeCondition struct {
	namedType namedType
	locator
}

func (t *typeCondition) NamedType() namedType {
	return t.namedType
}

type unionMemberTypes []*namedType

type implementsInterfaces []*namedType

type name struct {
	value string
	locator
}

func (n *name) Value() string {
	return n.value
}

type alias struct {
	name
}

type value interface {
	locatorInterface
	value() value
	Value() interface{}
}

var _ value = (*objectField)(nil)
var _ value = (*objectValue)(nil)
var _ value = (*listValue)(nil)
var _ value = (*intValue)(nil)
var _ value = (*floatValue)(nil)
var _ value = (*stringValue)(nil)
var _ value = (*booleanValue)(nil)
var _ value = (*enumValue)(nil)
var _ value = (*variable)(nil)

type objectField struct {
	name   name
	_value value
	locator
}

func (d *objectField) value() value {
	return d
}

func (of *objectField) Value() interface{} {
	return of._value.Value()
}

func (of *objectField) Name() name {
	return of.name
}

type objectValue struct {
	values []objectField
	locator
}

func (d *objectValue) value() value {
	return d
}

func (o *objectValue) Values() []objectField {
	return o.values
}

func (ov objectValue) Value() interface{} {
	return ov
}

type listValue struct {
	values []value
	locator
}

func (d *listValue) value() value {
	return d
}

func (lv listValue) Value() interface{} {
	return lv
}

func (l *listValue) Values() []value {
	return l.values
}

type intValue struct {
	_value int64
	locator
}

func (d *intValue) value() value {
	return d
}

func (iv *intValue) Value() interface{} {
	return iv._value
}

type floatValue struct {
	_value float64
	locator
}

func (d *floatValue) value() value {
	return d
}

func (fv *floatValue) Value() interface{} {
	return fv._value
}

type stringValue struct {
	_value string
	locator
}

func (d *stringValue) value() value {
	return d
}

func (sv *stringValue) Value() interface{} {
	return sv._value
}

type booleanValue struct {
	_value bool
	locator
}

func (d *booleanValue) value() value {
	return d
}

func (bv *booleanValue) Value() interface{} {
	return bv._value
}

type enumValue struct {
	name name
	locator
}

func (d *enumValue) value() value {
	return d
}

func (ev *enumValue) Value() interface{} {
	return ev.name
}

type nullValue struct {
	locator
}

func (d *nullValue) value() value {
	return d
}

func (nv *nullValue) Value() interface{} {
	return nil
}

type variable struct {
	name name
	locator
}

func (d *variable) value() value {
	return d
}

// TODO fetch variable value from variable map
func (v *variable) Value() interface{} {
	panic("not implemented")
	return nil
}

func (v *variable) Name() name {
	return v.name
}

type defaultValue struct {
	value value
	locator
}

func (dv defaultValue) Value() interface{} {
	return dv.value.Value()
}
