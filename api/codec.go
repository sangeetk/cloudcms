package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// encodeRequest encodes the request as JSON
func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeResponse decodes the response from the service
func decodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

// decodeSearchResults decodes search results
func decodeSearchResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results SearchResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// decodeFacetsSearchResults decodes facets search results
func decodeFacetsSearchResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results FacetsSearchResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// decodeListResults decodes list results
func decodeListResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results ListResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// decodeSchemaResponse decodes the response from the schema request
func decodeSchemaResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response SchemaResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
