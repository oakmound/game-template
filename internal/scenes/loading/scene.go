package loading

import (
	"path/filepath"
	"time"

	"github.com/disintegration/gift"
	"github.com/oakmound/game-template/internal/layers"
	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/grove/components/sound"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
)

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
		win := ctx.Window.(*oak.Window)
		// CONSIDER: Provide an ident to show an image that represents you as a creator
		ident, err := render.LoadSprite(filepath.Join("assets", "images", "raw", "ident.png"))
		if err == nil {
			ident.Modify(mod.ResizeToFit(400, 400, gift.CubicResampling))
			identw, identh := ident.GetDims()
			ident.SetPos(
				float64(ctx.Window.Width())/2-float64(identw)/2,
				float64(ctx.Window.Height())/2-float64(identh)/2,
			)
			ctx.DrawStack.Draw(ident, layers.StackBackground, layers.Back)
		}
		// CONSIDER: Providing some image that reassures users that yes the sytstem is doing work.
		loadSheet, err := render.LoadSheet(filepath.Join("assets", "images", "32x32", "loading.png"), intgeom.Point2{32, 32})
		if err == nil {
			loadingSeq, err = render.NewSheetSequence(loadSheet, 32, 0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0)
			loadingSeq.SetPos(float64(ctx.Window.Width()/2)-16, float64(ctx.Window.Height())-64)
			dlog.ErrorCheck(err)

			win.LoadingR = loadingSeq

			go ctx.DoAfter(2*time.Second, func() {
				// after some time, we start displaying the loading circle to reassure
				// the player things are still happening
				ctx.DrawStack.Draw(loadingSeq, layers.StackBackground, layers.Back)
			})
		}

		// Preload everything in the images folder
		go func() {
			imageFolder := "assets/images"
			if FastLoad {
				dlog.ErrorCheck(render.BlankBatchLoad(imageFolder, 1000*500))
			} else {
				dlog.ErrorCheck(render.BatchLoad(imageFolder))
			}
			loadComplete++
		}()

		// Load everything in the sound folder
		// You might consider initializing all your sound files here as well
		go func() {
			if FastLoad {
				audio.BlankBatchLoad(filepath.Join("assets", "audio"))
			} else {
				audio.BatchLoad(filepath.Join("assets", "audio"))
			}

			// CONSIDER: having an sfx package with the list of files to load with give sound fonts.

			sound.Init(1, 1, 1)
			loadComplete++
		}()
		go func() {
			waitInProduction(ctx)
			loadComplete++
		}()
	},
	Loop: func() bool {
		return loadComplete < 3
	},
	End: func() (string, *scene.Result) {
		return scenes.Sample, &scene.Result{
			Transition:     scene.Fade(1, 20),
			NextSceneInput: nil,
		}
	},
}
