package api

import (
	"github.com/pixality-inc/golang-boilerplate-project/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type ProtoRenderer struct{}

func NewProtoRenderer() *ProtoRenderer {
	return &ProtoRenderer{}
}

func (r *ProtoRenderer) Ok() proto.Message {
	return &protocol.OkResponse{}
}

func (r *ProtoRenderer) Error(_ int, err error) proto.Message {
	return &protocol.ErrorResponse{
		Error: &protocol.Error{
			Message: err.Error(),
		},
	}
}
