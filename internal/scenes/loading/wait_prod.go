//go:build prod
// +build prod

package loading

import (
	"time"

	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/scene"
)

func waitInProduction(ctx *scene.Context) {
	//press esc or wait 2 seconds
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.EventHandler.WaitForEvent(key.Down + key.Escape):
	}
}
