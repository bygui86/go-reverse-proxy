package blogpost

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-reverse-proxy/rest-app/logging"
)

const (
	errorMessageFormat         = "Error on %s (%s): %s"
	errorEncodeResponseMessage = "encode response"
	errorDecodeResponseMessage = "decode request"
	errorPostNotFoundMessage   = "BlogPost not found"
)

func (s *Server) getBlogPosts(writer http.ResponseWriter, request *http.Request) {
	requestStr := "GET BlogPosts request"
	logging.Log.Info(requestStr)

	setJsonContentType(writer)
	err := json.NewEncoder(writer).Encode(s.blogPosts)
	if err != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorEncodeResponseMessage, err.Error())
		setStatusInternalServerError(writer)
		returnErrorResponse(writer, requestStr, err.Error())
	}
}

func (s *Server) getBlogPostByPath(writer http.ResponseWriter, request *http.Request) {
	requestStr := "GET BlogPost request by path-var"
	logging.Log.Info(requestStr)

	pathVars := mux.Vars(request)
	if pathVars[idKey] == "" {
		errMsg := fmt.Sprintf("%s not found in URL", idKey)
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "retrieve path var", errMsg)
		setStatusBadRequest(writer)
		returnErrorResponse(writer, requestStr, errMsg)
		return
	}

	logging.SugaredLog.Infof("Searching for BlogPost by id %s", pathVars[idKey])
	for _, item := range s.blogPosts {
		if item.ID == pathVars[idKey] {
			setJsonContentType(writer)
			err := json.NewEncoder(writer).Encode(item)
			if err != nil {
				logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorEncodeResponseMessage, err.Error())
				setStatusInternalServerError(writer)
				returnErrorResponse(writer, requestStr, err.Error())
				return
			}
			return
		}
	}

	logging.SugaredLog.Warnf("Post with id %s not found", pathVars[idKey])
	setStatusNotFound(writer)
	returnErrorResponse(writer, requestStr, errorPostNotFoundMessage)
}

func (s *Server) getBlogPostByQuery(writer http.ResponseWriter, request *http.Request) {
	requestStr := "GET BlogPost by query request"
	logging.Log.Info(requestStr)

	id, queryErr := s.getIdFromQueryParams(request)
	if queryErr != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "retrieve query param", queryErr.Error())
		setStatusBadRequest(writer)
		returnErrorResponse(writer, requestStr, queryErr.Error())
		return
	}

	logging.SugaredLog.Infof("Searching for BlogPost by id %s", id)
	for _, item := range s.blogPosts {
		if item.ID == id {
			setJsonContentType(writer)
			encodeErr := json.NewEncoder(writer).Encode(item)
			if encodeErr != nil {
				logging.SugaredLog.Errorf(errorMessageFormat, requestStr, encodeErr.Error())
				setStatusInternalServerError(writer)
				returnErrorResponse(writer, requestStr, encodeErr.Error())
				return
			}
			return
		}
	}

	logging.SugaredLog.Warnf("Post with id %s not found", id)
	setStatusNotFound(writer)
	returnErrorResponse(writer, requestStr, errorPostNotFoundMessage)
}

func (s *Server) createBlogPost(writer http.ResponseWriter, request *http.Request) {
	requestStr := "CREATE BlogPost request"
	logging.Log.Info(requestStr)

	var post blogPost
	decErr := json.NewDecoder(request.Body).Decode(&post)
	if decErr != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorDecodeResponseMessage, decErr.Error())
		setStatusBadRequest(writer)
		returnErrorResponse(writer, requestStr, decErr.Error())
		return
	}

	post.ID = strconv.Itoa(rand.Intn(1000000))
	s.blogPosts = append(s.blogPosts, &post)

	setJsonContentType(writer)
	setStatusCreated(writer)
	encErr := json.NewEncoder(writer).Encode(&post)
	if encErr != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorEncodeResponseMessage, encErr.Error())
		setStatusInternalServerError(writer)
		returnErrorResponse(writer, requestStr, encErr.Error())
	}
}

func (s *Server) updateBlogPost(writer http.ResponseWriter, request *http.Request) {
	requestStr := "UPDATE BlogPosts request"
	logging.Log.Info(requestStr)

	pathVars := mux.Vars(request)
	if pathVars[idKey] == "" {
		errMsg := fmt.Sprintf("%s not found in URL", idKey)
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "retrieve path var", errMsg)
		setStatusBadRequest(writer)
		returnErrorResponse(writer, requestStr, errMsg)
		return
	}

	for index, item := range s.blogPosts {
		if item.ID == pathVars[idKey] {
			/*
				blogPosts[:index] >> from the beginning to index position
				blogPosts[index+1:]... >> from index+1 position to the end
			*/
			s.blogPosts = append(s.blogPosts[:index], s.blogPosts[index+1:]...)
			var post blogPost
			decErr := json.NewDecoder(request.Body).Decode(&post)
			if decErr != nil {
				logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorDecodeResponseMessage, decErr.Error())
				setStatusBadRequest(writer)
				returnErrorResponse(writer, requestStr, decErr.Error())
				return
			}

			post.ID = pathVars[idKey]
			s.blogPosts = append(s.blogPosts, &post)

			setJsonContentType(writer)
			setStatusAccepted(writer)
			encErr := json.NewEncoder(writer).Encode(&post)
			if encErr != nil {
				logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorEncodeResponseMessage, encErr.Error())
				setStatusInternalServerError(writer)
				returnErrorResponse(writer, requestStr, encErr.Error())
				return
			}
			return
		}
	}

	setStatusNotFound(writer)
	returnErrorResponse(writer, requestStr, errorPostNotFoundMessage)
}

func (s *Server) deleteBlogPost(writer http.ResponseWriter, request *http.Request) {
	requestStr := "DELETE BlogPosts request"
	logging.Log.Info(requestStr)

	pathVars := mux.Vars(request)
	if pathVars[idKey] == "" {
		errMsg := fmt.Sprintf("%s not found in URL", idKey)
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "retrieve path var", errMsg)
		setStatusBadRequest(writer)
		returnErrorResponse(writer, requestStr, errMsg)
		return
	}

	for index, item := range s.blogPosts {
		if item.ID == pathVars[idKey] {
			s.blogPosts = append(s.blogPosts[:index], s.blogPosts[index+1:]...)
			break
		}
	}

	setStatusAccepted(writer)
}

func (s *Server) getRoot(writer http.ResponseWriter, request *http.Request) {
	requestStr := "GET root request"
	logging.Log.Info(requestStr)

	_, err := writer.Write([]byte(rootMessage))
	if err != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "write string response", err.Error())
		setStatusInternalServerError(writer)
		returnErrorResponse(writer, requestStr, err.Error())
	}
}

func (s *Server) getRoutes(writer http.ResponseWriter, request *http.Request) {
	requestStr := "GET routes request"
	logging.Log.Info(requestStr)

	if s.routes == nil || len(s.routes) < 1 {
		walkErr := s.router.Walk(s.walkRoute)
		if walkErr != nil {
			logging.SugaredLog.Errorf(errorMessageFormat, requestStr, "walk through routes", walkErr.Error())
			setStatusInternalServerError(writer)
			returnErrorResponse(writer, requestStr, walkErr.Error())
			return
		}
	}

	setJsonContentType(writer)
	err := json.NewEncoder(writer).Encode(s.routes)
	if err != nil {
		logging.SugaredLog.Errorf(errorMessageFormat, requestStr, errorEncodeResponseMessage, err.Error())
		setStatusInternalServerError(writer)
		returnErrorResponse(writer, requestStr, err.Error())
	}
}
