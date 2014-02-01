// Copyright (c) 2014 Jason Goecke
// client.go

package wit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	USERAGENT = "WIT (Go net/http)"
)

// Represents a client for the Wit API (https://wit.ai/docs/api)
type WitClient struct {
	ApiBase string
}

// Represents the HTTP parameters to pass along to the Wit API
type HttpParams struct {
	Verb        string
	Resource    string
	ContentType string
	Data        []byte
}

// Stores the ApiKey for the Wit API
var ApiKey string

// Creates a NewClient for the Wit API
//
//		client := wit.NewClient("<ACCESS-TOKEN>")
func NewClient(apiKey string) *WitClient {
	client := &WitClient{}
	client.ApiBase = "https://api.wit.ai"
	ApiKey = apiKey
	return client
}

// Provides a common facility for doing a DELETE on a Wit resource
//
//		result, err := delete("https://api.wit.ai/entities", "favorite_city")
func delete(resource string, id string) ([]byte, error) {
	httpParams := &HttpParams{}
	httpParams.Resource = resource + "/" + id
	httpParams.Verb = "DELETE"
	return processRequest(httpParams)
}

// Provides a common facility for doing a GET on a Wit resource
//
//		result, err := get("https://api.wit.ai/entities/favorite_city")
func get(resource string) ([]byte, error) {
	httpParams := &HttpParams{}
	httpParams.Resource = resource
	httpParams.Verb = "GET"
	return processRequest(httpParams)
}

// Provides a common facility for doing a POST on a Wit resource. Takes
// JSON []byte for the data argument.
//
//		result, err := post("https://api.wit.ai/entities", entity)
func post(resource string, data []byte) ([]byte, error) {
	httpParams := &HttpParams{"POST", resource, "application/json", data}
	return processRequest(httpParams)
}

// Provides a common facility for doing a POST with a file on a Wit resource.
//
//		result, err := postFile("https://api.wit.ai/messages", message)
func postFile(resource string, request *MessageRequest) ([]byte, error) {
	if request.File != "" {
		file, err := os.Open(request.File)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		stats, statsErr := file.Stat()
		if statsErr != nil {
			return nil, statsErr
		}
		var size int64 = stats.Size()
		data := make([]byte, size)
		file.Read(data)
		httpParams := &HttpParams{"POST", resource, request.ContentType, data}
		return processRequest(httpParams)
	} else {
		if request.FileContents != nil {
			httpParams := &HttpParams{"POST", resource, request.ContentType, request.FileContents}
			return processRequest(httpParams)
		} else {
			return nil, errors.New("Must provide a filename or contents")
		}
	}
}

// Provides a common facility for doing a PUT on a Wit resource.
//
//		result, err := put("https://api.wit.ai/entities", entity)
func put(resource string, data []byte) ([]byte, error) {
	httpParams := &HttpParams{"PUT", resource, "application/json", data}
	return processRequest(httpParams)
}

// Processes an HTTP request to the Wit API
func processRequest(httpParams *HttpParams) ([]byte, error) {
	reader := bytes.NewReader(httpParams.Data)
	httpClient := &http.Client{}
	req, err := http.NewRequest(httpParams.Verb, httpParams.Resource, reader)
	if err != nil {
		return nil, err
	}
	setHeaders(req, httpParams.ContentType)
	result, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(result.Body)
	result.Body.Close()
	if result.StatusCode != 200 {
		return nil, errors.New(http.StatusText(result.StatusCode))
	}
	return body, nil
}

// Sets the custom headers required for the Wit.ai API
//
//		setHeaders(req, httpParams.ContentType)
func setHeaders(req *http.Request, contentType string) {
	req.Header.Add("Authorization", "Bearer "+ApiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
}
