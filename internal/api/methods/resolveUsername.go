package methods

import (
	"elide/pkg/elide"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anonyindian/gotgproto/ext"
)

func ResolveUsername(ctx *ext.Context, body json.RawMessage) (any, error) {
	var p elide.ResolveUsernameBody
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to read body for resolveUsername: %s", err.Error())
	}
	if p.Username == "" {
		return nil, errors.New("username not provided")
	}
	return ctx.ResolveUsername(p.Username)
}
