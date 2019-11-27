package caf

import (
	"github.com/omeryahud/caf/internal/pkg/validators/graphqlvalidator"
	"log"
	"net/http"
	"strconv"

	"github.com/omeryahud/caf/internal/pkg/configs"
	"github.com/omeryahud/caf/internal/pkg/middleman"
	"github.com/omeryahud/caf/internal/pkg/proxymiddlewares"
	"github.com/omeryahud/caf/internal/pkg/validators"
	"github.com/omeryahud/caf/internal/pkg/validators/jsonvalidator"
	"github.com/pkg/errors"
)

// AddValidationMiddlewares gets a reference to a Middleman and a slice of targets
// and creates a new middleware for each endpoint in the targets' apis.
func addValidationMiddlewares(mm *middleman.Middleman, targets []configs.Target) error {
	// Loop over the targets slice
	for _, target := range targets {
		// For each target loop over its apis
		for index, api := range target.Apis {
			var err error
			var validator validators.Validator

			// Each api has a validator that filter the api's traffic.
			// Here we decide which validator to create according to the api's type.
			switch api.Type {
			case configs.TypeRest:
				validator, err = jsonvalidator.NewJsonValidator(api.Version)
				if err != nil {
					return errors.Wrap(err, "failed to created validator for number - "+strconv.Itoa(index))
				}
			case configs.TypeGraphQL:
				validator, err = graphqlvalidator.NewGraphQLValidator()
				if err != nil {
					return errors.Wrap(err, "failed to create validator for api number - "+strconv.Itoa(index))
				}
			default:
				log.Print("[Proxy WARNING]: Invalid API Type - " + api.Type)
			}

			// For each api loop over its endpoints
			for _, endpoint := range api.Endpoints {
				//Add the endpoint's schema to the api's validator.
				err := validator.LoadSchema(endpoint.Path, endpoint.Method, []byte(endpoint.Schema))
				if err != nil {
					log.Print("[Proxy ERROR]: Failed to load schema for endpoint - " + endpoint.Path + ", Error: " + err.Error())
					return err
				}

				// Creating a new ValidateRequest middleware with the appropriate HTTP method.
				switch endpoint.Method {
				case http.MethodGet:
					mm.Get(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodPost:
					mm.Post(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodPut:
					mm.Put(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodDelete:
					mm.Delete(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case "ALL":
					mm.All(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				default:
					log.Print("[Proxy WARNING]: Invalid method - " + endpoint.Method + " for endpoint - " + endpoint.Path)
				}

				log.Print("[Proxy DEBUG]: Added middleware for - " + endpoint.Method + " " + endpoint.Path)
			}
		}
	}

	return nil
}
