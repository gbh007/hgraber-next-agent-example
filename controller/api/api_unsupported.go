package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-example/controller/api/internal/server"
)

func (c *Controller) APIFsCreatePost(ctx context.Context, req server.APIFsCreatePostReq, params server.APIFsCreatePostParams) (server.APIFsCreatePostRes, error) {
	return &server.APIFsCreatePostBadRequest{
		InnerCode: ValidationCode,
		Details:   server.NewOptString("unsupported api"),
	}, nil
}

func (c *Controller) APIFsDeletePost(ctx context.Context, req *server.APIFsDeletePostReq) (server.APIFsDeletePostRes, error) {
	return &server.APIFsDeletePostBadRequest{
		InnerCode: ValidationCode,
		Details:   server.NewOptString("unsupported api"),
	}, nil
}

func (c *Controller) APIFsGetGet(ctx context.Context, params server.APIFsGetGetParams) (server.APIFsGetGetRes, error) {
	return &server.APIFsGetGetBadRequest{
		InnerCode: ValidationCode,
		Details:   server.NewOptString("unsupported api"),
	}, nil
}

func (c *Controller) APIFsIdsGet(ctx context.Context) (server.APIFsIdsGetRes, error) {
	return &server.APIFsIdsGetBadRequest{
		InnerCode: ValidationCode,
		Details:   server.NewOptString("unsupported api"),
	}, nil
}
