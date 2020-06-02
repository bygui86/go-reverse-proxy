package blogpost

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-reverse-proxy/rest-app/logging"
)

func buildBlogPost(id, title, body string) *blogPost {
	return &blogPost{ID: id, Title: title, Body: body}
}

func buildErrorResponse(request string, errorMsg string) *errorResponse {
	return &errorResponse{Request: request, Message: errorMsg}
}

func returnErrorResponse(writer http.ResponseWriter, request, errorMsg string) {
	err := json.NewEncoder(writer).Encode(buildErrorResponse(request, errorMsg))
	if err != nil {
		logging.SugaredLog.Errorf("Error on %s (encode ERROR response): %s - No response back to client",
			request, err.Error())
	}
}

func (s *Server) getIdFromQueryParams(request *http.Request) (string, error) {
	queryParams := request.URL.Query()
	idValues, idValuesFound := queryParams[idKey]
	id := idValues[0]
	if !idValuesFound || len(id) < 1 {
		return "", fmt.Errorf("%s query param not found", idKey)
	}
	return id, nil
}

// INFO: no need to test this function, it should be already tested by "gorilla/mux" library
func (s *Server) walkRoute(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	serverRoute := &serverRoute{}
	pathTemplate, pathTempErr := route.GetPathTemplate()
	if pathTempErr != nil {
		return pathTempErr
	}
	pathRegexp, pathRegErr := route.GetPathRegexp()
	if pathRegErr != nil {
		return pathRegErr
	}
	queriesTemplates, queriesTempErr := route.GetQueriesTemplates()
	if queriesTempErr != nil {
		return queriesTempErr
	}
	queriesRegexps, queriesRegErr := route.GetQueriesRegexp()
	if queriesRegErr != nil {
		return queriesRegErr
	}
	methods, methodsErr := route.GetMethods()
	if methodsErr != nil {
		return methodsErr
	}

	serverRoute.PathTemplate = pathTemplate
	serverRoute.PathRegexp = pathRegexp
	serverRoute.QueriesTemplate = queriesTemplates
	serverRoute.QueriesRegexp = queriesRegexps
	serverRoute.Methods = methods

	s.routes = append(s.routes, serverRoute)
	return nil
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setJsonContentType(writer http.ResponseWriter) {
	writer.Header().Set(contentTypeHeaderKey, applicationJsonValue)
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setStatusCreated(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusCreated)
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setStatusAccepted(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusAccepted)
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setStatusNotFound(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNotFound)
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setStatusBadRequest(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusBadRequest)
}

// INFO: no need to test this function, it should be already tested by "net/http" library
func setStatusInternalServerError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
}
