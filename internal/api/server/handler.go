package server

import (
	"encoding/json"

	"github.com/anonyindian/gotgproto/ext"
)

type HandlerFunc func(ctx *ext.Context, body json.RawMessage) (any, error)
