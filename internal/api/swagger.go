package api

import (
	"github.com/damek86/url-shortener-go/internal/config"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"log"
	"net/http"
)

func SetupSwaggerUiServing(wsContainer *restful.Container, cfg config.StaticContent) {
	if cfg.Disabled {
		log.Print("Serving Swagger-UI disabled!")
	} else {
		log.Printf("Serving Swagger-UI at %s", cfg.HttpPath)
		wsContainer.Handle(cfg.HttpPath,
			http.StripPrefix(cfg.HttpPath,
				http.FileServer(http.Dir(cfg.FilePath))))
	}
}

func CreateSwaggerConfig(services []*restful.WebService) restfulspec.Config {
	return restfulspec.Config{
		WebServices: services,
		APIPath:     "/apidocs.json",
		DisableCORS: true,
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			injectDefinitionExamples(s)
		},
	}
}

var exampleDefinitions = map[string]interface{}{}

// NOTE: Unfortunately the 'go-openapi' extension of 'go-restful' does not provide support for "example" values within the OpenApi-Specification.
// As workaround we use a custom solution of "late injection" of examples. (it was inspired by https://swdc.visualstudio.com/SW%20Framework/_git/common-restfulspecext)
func injectDefinitionExamples(swagger *spec.Swagger) {
	for exampleType, example := range exampleDefinitions {
		if schema, ok := swagger.Definitions[exampleType]; ok {
			schema.Example = example
			swagger.Definitions[exampleType] = schema
		}
	}
}
