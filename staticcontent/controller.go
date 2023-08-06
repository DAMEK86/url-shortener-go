package staticcontent

import (
	"github.com/damek86/url-shortener-go/config"
	"github.com/emicklei/go-restful/v3"
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
