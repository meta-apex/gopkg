package rescue

import (
	"context"
	"fmt"

	"github.com/meta-apex/gopkg/zlog"
)

// Recover is used with defer to do cleanup on panics.
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		zlog.Error().Stack().Msg(fmt.Sprint(p))
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		zlog.Error().Stack().Msg(fmt.Sprint(p))
	}
}
