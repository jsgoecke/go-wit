// Copyright (c) 2014 Jason Goecke
// client.go

package wit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
)

const (
	// UserAgent is the HTTP Uesr Agent sent on HTTP requests
	UserAgent = "WIT (Go net/http)"
	// APIVersion is the version of the Wit API supported
	APIVersion = "v=20151127"
)

// Client represents a client for the Wit API (https://wit.ai/docs/api)
type Client struct {
	APIBase string
}

// HTTPParams represents the HTTP parameters to pass along to the Wit API
type HTTPParams struct {
	Verb        string
	Resource    string
	ContentType string
	Data        []byte
}

// Stores the ApiKey for the Wit API
var APIKey string

// NewClient creates a new client for the Wit API
//
//		client := wit.NewClient("<ACCESS-TOKEN>")
func NewClient(apiKey string) *Client {
	client := &Client{APIBase: "https://api.wit.ai"}
	APIKey = apiKey
	return client
}

// Provides a common facility for doing a DELETE on a Wit resource
//
//		result, err := delete("https://api.wit.ai/entities", "favorite_city")
func delete(resource string, id string) ([]byte, error) {
	httpParams := &HTTPParams{
		Resource: resource + "/" + id,
		Verb:     "DELETE",
	}
	return processRequest(httpParams)
}

// Provides a common facility for doing a GET on a Wit resource
//
//		result, err := get("https://api.wit.ai/entities/favorite_city")
func get(resource string) ([]byte, error) {
	httpParams := &HTTPParams{
		Resource: resource,
		Verb:     "GET",
	}
	return processRequest(httpParams)
}

// Provides a common facility for doing a POST on a Wit resource. Takes
// JSON []byte for the data argument.
//
//		result, err := post("https://api.wit.ai/entities", entity)
func post(resource string, data []byte) ([]byte, error) {
	httpParams := &HTTPParams{"POST", resource, "application/json", data}
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
		size := stats.Size()
		data := make([]byte, size)
		file.Read(data)
		httpParams := &HTTPParams{"POST", resource, request.ContentType, data}
		return processRequest(httpParams)
	}

	if request.FileContents != nil {
		httpParams := &HTTPParams{"POST", resource, request.ContentType, request.FileContents}
		return processRequest(httpParams)
		// } else {
		// return nil, errors.New("Must provide a filename or contents")
	}

	return nil, errors.New("must provide a filename or contents")
}

// Provides a common facility for doing a PUT on a Wit resource.
//
//		result, err := put("https://api.wit.ai/entities", entity)
func put(resource string, data []byte) ([]byte, error) {
	httpParams := &HTTPParams{"PUT", resource, "application/json", data}
	return processRequest(httpParams)
}

// Processes an HTTP request to the Wit API
func processRequest(httpParams *HTTPParams) ([]byte, error) {
	regex := regexp.MustCompile(`\?`)
	if regex.MatchString(httpParams.Resource) {
		httpParams.Resource += "&" + APIVersion
	} else {
		httpParams.Resource += "?" + APIVersion
	}
	reader := bytes.NewReader(httpParams.Data)
	httpClient := &http.Client{}
	req, err := http.NewRequest(httpParams.Verb, httpParams.Resource, reader)
	if err != nil {
		return nil, err
	}
	setHeaders(req, httpParams.ContentType)

	if os.Getenv("GOWIT_DEBUG") == "true" {
		debug(httputil.DumpRequestOut(req, true))
	}

	result, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if os.Getenv("GOWIT_DEBUG") == "true" {
		debug(httputil.DumpResponse(result, true))
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
	req.Header.Add("Authorization", "Bearer "+APIKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
}

func debug(data []byte, err error) {
	if err == nil {
		if len(data) > 1000 {
			fmt.Printf("DATA TOO LARGE %d\n\n", len(data))
		} else {
			fmt.Printf("%s\n\n", data)
		}
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
