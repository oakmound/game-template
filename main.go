package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"time"

	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/game-template/internal/scenes/debug"
	"github.com/oakmound/game-template/internal/scenes/loading"
	"github.com/oakmound/game-template/internal/scenes/sample"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
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
			c.EventRefreshRate = oak.Duration(100 * time.Millisecond)
			c.EnableDebugConsole = debugEnabled
			c.BatchLoad = false

			return c, nil
		})
		return fmt.Errorf("finished main window with output of: %v", err)
	})

	errG.Wait()
}
