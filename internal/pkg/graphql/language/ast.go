package language

type DirectiveLocation int

const (
	// Executable Directive Locations
	QUERY DirectiveLocation = iota + 1
	MUTATION
	SUBSCRIPTION
	FIELD
	FRAGMENT_DEFINITION
	FRAGMENT_SPREAD
	INLINE_FRAGMENT
	VARIABLE_DEFINITION

	// Type System Directive Locations
	SCHEMA
	SCALAR
	OBJECT
	FIELD_DEFINITION
	ARGUMENT_DEFINITION
	INTERFACE
	UNION
	ENUM
	ENUM_VALUE
	INPUT_OBJECT
	INPUT_FIELD_DEFINITION
)

type TypeSystemDirectiveLocation DirectiveLocation

type ExecutableDirectiveLocation DirectiveLocation

type DirectiveLocations []DirectiveLocation

type name string

func parseName(n string) name {
	return name("")
}

type DirectiveDefinition struct {
	Description         *string
	Name                name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
}

type ArgumentsDefinition []InputValueDefinition

type TypeKind int

const (
	NAMED_TYPE TypeKind = iota + 1
	LIST_TYPE
	NON_NULL_TYPE
)

type Type interface {
	GetTypeKind() TypeKind
}

type NamedType name

func (nt NamedType) GetTypeKind() TypeKind {
	return NAMED_TYPE
}

type ListType struct {
	Type Type
}

func (lt ListType) GetTypeKind() TypeKind {
	return LIST_TYPE
}

type NonNullType struct {
	Type Type
}

func (nt NonNullType) GetTypeKind() TypeKind {
	return NON_NULL_TYPE
}
