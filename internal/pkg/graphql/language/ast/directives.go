package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

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

type Directives []Directive

type Directive struct {
	Name      Name
	Arguments *Arguments
	Loc       *location.Location
}
