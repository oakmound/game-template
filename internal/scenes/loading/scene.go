package loading

import (
	"time"

	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

//PreLoadTimeStr
const PreLoadTimeStr = "preloadtime"

var (
	loadComplete = 0
	loadingSeq   render.Renderable
)

// FastLoad skips the pesky loading of images and just makes a bunch of empty images
// When your asset library is large and you just want to test some basic logic.
var FastLoad bool

// Scene for managing loading and displaying something while we load.
var Scene = scene.Scene{
	Start: func(ctx *scene.Context) {
		// Fake a longer wait time by ensuring that in production we show the loading information
		// such as a loading splash.
		go func() {
			if t := ctx.Value(PreLoadTimeStr); t != nil {
				if preTime, ok := t.(time.Time); ok {
					waitInProduction(ctx, preTime)
				}
			}
			ctx.Window.NextScene()
		}()

	},
	End: func() (string, *scene.Result) {
		return scenes.Sample, &scene.Result{
			Transition:     scene.Fade(1, 20),
			NextSceneInput: nil,
		}
	},
}
