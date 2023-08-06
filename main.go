package main

import (
	"context"
	"fmt"
	"github.com/damek86/url-shortener-go/api"
	"github.com/damek86/url-shortener-go/config"
	"github.com/damek86/url-shortener-go/staticcontent"
	"github.com/damek86/url-shortener-go/url"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.DefaultConfig

	router := new(restful.WebService)

	router.Path("/api")
	api.RegisterRoutes(router, api.Services{
		UrlShortener: url.NewService(),
	})

	container := restful.NewContainer()
	// Assume all requests to be authorized using the Authorization header, so this is safe
	container.Filter(restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE", "PATCH"},
		ExposeHeaders:  []string{"Location"},
		CookiesAllowed: true, // needed to allow credentials
	}.Filter)

	container.Add(router)

	swaggerConfig := api.CreateSwaggerConfig(container.RegisteredWebServices())
	openApiService := restfulspec.NewOpenAPIService(swaggerConfig)
	container.Add(openApiService)

	staticcontent.SetupSwaggerUiServing(container, cfg.Swagger)

	defaultTimeout := 10 * time.Second
	port := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: container.ServeMux,
		//ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		IdleTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
	}
	// everything below enables graceful shutdown without dropping any requests
	go func() {
		log.Printf("Listening on port %d", port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen on port %d. Shutting down!", port)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("Shutting down server...")
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Printf("Could not shutdown server cleanly")
		os.Exit(1)
	}
}
