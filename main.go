package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/disintegration/gift"
	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/game-template/internal/scenes/debug"
	"github.com/oakmound/game-template/internal/scenes/loading"
	"github.com/oakmound/game-template/internal/scenes/sample"
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"golang.org/x/sync/errgroup"
)

//go:embed assets
var assets embed.FS

var debugFlag = flag.Bool("debug", false, "turn on debug")

func main() {
	// Load flags
	flag.Parse()

	// This is not a great method but shows how one can pass flags into the initial oak call for operational updates.
	debugEnabled := *debugFlag

	oak.SetFS(assets)

	// A common set of draw stacks, Composite as background, dynamic as foreground and then the static layer as a UI type thing.
	// See internal/layering/layers.go for the names.
	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)

	errG, ctx := errgroup.WithContext(context.Background())

	// Add scenes from local internal packages
	mainWindow := oak.NewWindow()
	mainWindow.ParentContext = context.WithValue(ctx, "name", "Sample App")
	mainWindow.AddScene(scenes.Batchloading, loading.Scene)
	mainWindow.AddScene(scenes.Sample, sample.Scene)

	if debugEnabled {

		dlog.Info("Enabled Debug mode and console")

		debugger := oak.NewWindow()
		debugger.ParentContext = context.WithValue(mainWindow.ParentContext, "name", "debugger")
		debugger.ParentContext = context.WithValue(debugger.ParentContext, debug.ExtraBus, mainWindow.EventHandler())
		fmt.Println(debugger.AddScene(scenes.Debug, debug.Scene))

		secondaryCaller := event.NewCallerMap()
		debugger.SetLogicHandler(event.NewBus(secondaryCaller))
		debugger.CallerMap = secondaryCaller
		debugger.DrawStack = render.NewDrawStack(render.NewStaticHeap())
		errG.Go(func() error {
			err := debugger.Init(scenes.Debug, func(c oak.Config) (oak.Config, error) {
				c.Debug.Level = "VERBOSE"
				c.Title = "Debugging Console"
				c.Screen = oak.Screen{
					Width:  960,
					Height: 360,
					Scale:  1,
				}
				c.TopMost = true
				c.FrameRate = 30
				c.DrawFrameRate = 120
				c.EnableDebugConsole = true
				return c, nil
			})
			return fmt.Errorf("finished debug window with output of: %v", err)
		})
	}

	mainWindow.ParentContext = context.WithValue(mainWindow.ParentContext, loading.PreLoadTimeStr, time.Now())

	screenR := render.NewColorBox(mainWindow.Bounds().X(), mainWindow.Bounds().Y(), color.RGBA{255, 255, 22, 0})
	mid := mainWindow.Bounds().DivConst(2)
	// CONSIDER: Providing some image that reassures users that yes the sytstem is doing work.
	// loadSheet, err := render.LoadSheet(filepath.Join("assets", "images", "32x32", "loading.png"), intgeom.Point2{32, 32})
	// if err == nil {
	// 	loadingSeq, err := render.NewSheetSequence(loadSheet, 32, 0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0)
	// 	// loadingSeq.SetPos(float64(ctx.Window.Bounds().X()/2)-16, float64(ctx.Window.Bounds().Y())-64)
	// 	dlog.ErrorCheck(err)

	// 	// go ctx.DoAfter(2*time.Second, func() {
	// 	// 	// after some time, we start displaying the loading circle to reassure
	// 	// 	// the player things are still happening
	// 	// 	ctx.DrawStack.Draw(loadingSeq, layers.StackBackground, layers.Back)
	// 	// })
	// } else {
	// 	fmt.Prindt have loading.png", err.Error())
	// }
	// CONSIDER: Provide an ident to show an image that represents you as a creator
	ident, err := render.LoadSprite(filepath.Join("assets", "images", "raw", "ident.png"))
	if err == nil {
		ident.Modify(mod.ResizeToFit(400, 400, gift.CubicResampling))
		// identw, identh := ident.GetDims()
		// ident.SetPos(
		// 	float64().Bounds().X())/2-float64(identw)/2,
		// 	float64(ctx.Window.Bounds().Y())/2-float64(identh)/2,
		// )
		// ctx.DrawStack.Draw(ident, layers.StackBackground, layers.Back)
		ident.Draw(screenR, float64(mid.X()), float64(mid.Y()))
	}
	mainWindow.LoadingR = screenR

	errG.Go(func() error {
		err := mainWindow.Init(scenes.Batchloading, func(c oak.Config) (oak.Config, error) {
			// Setup oak config

			c.FrameRate = 60
			c.DrawFrameRate = 60
			c.Screen = oak.Screen{
				Width:  1280,
				Height: 720,
				Scale:  1,
			}
			c.Debug = oak.Debug{
				Level: "Info",
			}
			c.Title = "Sample App"
			c.TrackInputChanges = true
			c.LoadBuiltinCommands = true
			c.TopMost = true
			c.EnableDebugConsole = debugEnabled
			c.BatchLoad = false

			return c, nil
		})
		return fmt.Errorf("finished main window with output of: %v", err)
	})

	errG.Wait()
}
