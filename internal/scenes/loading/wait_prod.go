//go:build prod
// +build prod

package loading

import (
	"time"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/scene"
)

func waitInProduction(ctx *scene.Context) {

	escChan := make(chan struct{})
	event.GlobalBind(ctx, key.Down(key.Escape), func(key.Event) event.Response {
		escChan <- struct{}{}
		return 1
	})

	//press esc or wait 2 seconds
	select {
	case <-time.After(2 * time.Second):
	case <-escChan:
	}
}
