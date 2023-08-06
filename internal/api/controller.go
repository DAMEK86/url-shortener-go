package api

import (
	"encoding/json"
	"fmt"
	"github.com/damek86/url-shortener-go/internal/config"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	buildVersionHeader  = "X-BuildVersion"
	sourceVersionHeader = "X-SourceVersion"
)

type controller struct {
	services Services
}

func (c *controller) handleGetHealth(_ *restful.Request, resp *restful.Response) {
	resp.ResponseWriter.Header().Set(buildVersionHeader, config.BuildVersion)
	resp.ResponseWriter.Header().Set(sourceVersionHeader, config.SourceVersion)
}

func (c *controller) handleShortUrl(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	urlRequest := UrlRequest{}
	if err := GetJsonBody(req.Request, &urlRequest); err != nil {
		WriteErrorResponse(resp, http.StatusBadRequest, InvalidParameterErrorID, err.Error())
		return
	}

	shortenUrl, err := c.services.UrlShortener.ShortenUrl(ctx, urlRequest.Url)
	if err != nil {
		WriteErrorResponse(resp.ResponseWriter, http.StatusBadRequest, InvalidParameterErrorID, fmt.Sprintf("Invalid 'URL' value: %s", err.Error()))
		return
	}
	WriteJsonResult(resp.ResponseWriter, UrlResponse{
		Url: shortenUrl,
	})
}

func (c *controller) handleResolveUrl(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	shortenUrl, err := GetPathParam(req, "id")
	if err != nil {
		WriteErrorResponse(resp.ResponseWriter, http.StatusBadRequest, InvalidParameterErrorID, fmt.Sprintf("Invalid 'id' value: %s", err.Error()))
		return
	}
	resolvedUrl, err := c.services.UrlShortener.ResolveUrl(ctx, shortenUrl)
	if err != nil {
		WriteErrorResponse(resp.ResponseWriter, http.StatusNotFound, EntityNotFoundErrorID, fmt.Sprintf("No Url found for id: %s", err.Error()))
		return
	}

	http.Redirect(resp.ResponseWriter, req.Request, resolvedUrl, http.StatusSeeOther)
}

func WriteJsonResult(writer http.ResponseWriter, result interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		r := &ErrorResponse{
			ErrorID:      UnexpectedErrorID,
			ErrorMessage: err.Error(),
		}
		err = encoder.Encode(r)
		if err != nil {
			log.Print("failed to serialize error response")
		}
	}
}

func GetPathParam(request *restful.Request, name string) (string, error) {
	pathValue := GetFromPath(request, name)
	if len(strings.TrimSpace(pathValue)) == 0 {
		return "", fmt.Errorf("Parameter '%s' is missing", name)
	}
	return pathValue, nil
}

func GetFromPath(request *restful.Request, pathParam string) string {
	value, _ := url.PathUnescape(request.PathParameter(pathParam))
	return value
}

func WriteErrorResponse(writer http.ResponseWriter, httpStatusCode int, errorID string, errorMsg string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)
	r := &ErrorResponse{
		ErrorID:      errorID,
		ErrorMessage: errorMsg,
	}
	WriteJsonResult(writer, r)
}

func GetJsonBody(request *http.Request, data interface{}) error {
	if request.Body == nil || request.ContentLength == 0 {
		return fmt.Errorf("Missing JSON body")
	}
	err := json.NewDecoder(request.Body).Decode(data)
	if err != nil {
		return fmt.Errorf("Error while deserializing JSON body: %w", err)
	}
	return nil
}
