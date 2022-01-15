package debug

import (
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/scene"
)

// RecreationNeeded tracks when we should recheck for main loop changes
const RecreationNeeded = "RecreationNeeded"

type phantomEntity struct {
	event.CID
	outsideCID event.CID
}

// Init the phantomEntity on the main event handler so that we can perform future checks
func (pe *phantomEntity) Init() event.CID {
	return event.NextID(pe)
}

// start a watchdog that checks for status of the default event stream
// and publishes a recreationneeded event for others on its current scenes event handler
// This is really lossy and will leak events at starts of scenes.
func startWatchdog(ctx *scene.Context) {

	tester := &phantomEntity{}
	tester.CID = ctx.CallerMap.NextID(tester)

	tester.outsideCID = tester.Init()

	// Register a handler to the watchdog for our own verbose logging for recreation events.
	ctx.EventHandler.GlobalBind(RecreationNeeded, func(id event.CID, payload interface{}) int {
		ctx.EventHandler.Trigger("DisplayText", "Main Bus was reset")
		// ctx.Window.ShowNotification("reset", "Recieved a recreationRequest", false)
		return 0
	})

	// Start the watchdog go func rely on the ctx's cancellation for cleanup
	go func() {
		for {
			ctx.DoAfter(time.Second, func() {
				var stillOK bool
				testa := event.DefaultCallerMap.GetEntity(tester.outsideCID)
				if testa != nil {
					_, stillOK = testa.(*phantomEntity)
				}
				if stillOK {
					return
				}
				tester.outsideCID = tester.Init()
				ctx.EventHandler.Trigger("RecreationNeeded", struct{}{})

			})
		}
	}()

}
