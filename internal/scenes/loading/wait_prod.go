//go:build prod
// +build prod

package loading

import (
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/scene"
	"time"
)

func waitInProduction(ctx *scene.Context, start time.Time) {
	//press esc or wait 2 seconds
	delayTime := time.Until(start.Add(10 * time.Second))
	cancelCh := make(chan struct{})
	event.GlobalBind(ctx, key.Down(key.Escape), func(ev key.Event) event.Response {
		cancelCh <- struct{}{}
		return 1
	})
	select {
	case <-time.After(delayTime):
	case <-cancelCh:
	}
}
