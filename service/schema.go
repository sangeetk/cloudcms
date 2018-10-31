package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/item"
	"github.com/urantiatech/kit/endpoint"
)

// Schema - explains the schema
func (s *Service) Schema(ctx context.Context, req *api.SchemaRequest) (*api.SchemaResponse, error) {
	var resp = api.SchemaResponse{Schema: make(map[string]api.ContentType)}

	for t, v := range item.Fields {
		contentType := api.ContentType{
			Fields: v,
		}
		resp.Schema[t] = contentType
	}

	return &resp, nil
}

// SchemaEndpoint - creates endpoint for Schema service
func SchemaEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.SchemaRequest)
		return svc.Schema(ctx, &req)
	}
}

// DecodeSchemaReq - decodes the incoming request
func DecodeSchemaReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.SchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
