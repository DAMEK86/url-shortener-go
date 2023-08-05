package api

import (
	"github.com/damek86/url-shortner-go/url"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

type Services struct {
	UrlShortener url.Service
}

func RegisterRoutes(router *restful.WebService, services Services) {
	c := &controller{
		services: services,
	}

	router.Route(
		router.GET("/health").
			To(c.handleGetHealth).
			Doc("Gets health information about the service.").
			Metadata(restfulspec.KeyOpenAPITags, []string{"info"}).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)

	router.Route(
		router.POST("admin/u").
			To(c.handleShortUrl).
			Doc("Shorten provided url").
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(UrlRequest{}).
			Metadata(restfulspec.KeyOpenAPITags, []string{"url"}).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), UrlResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}),
	)

	router.Route(
		router.GET("u/{id}").
			To(c.handleResolveUrl).
			Param(restful.PathParameter("id", "shorten url").
				DataType("string").Required(true)).
			Doc("Shorten provided url").
			Metadata(restfulspec.KeyOpenAPITags, []string{"url"}).
			Returns(http.StatusSeeOther, http.StatusText(http.StatusSeeOther), nil).
			Returns(http.StatusNotFound, http.StatusText(http.StatusInternalServerError), ErrorResponse{}),
	)
}
