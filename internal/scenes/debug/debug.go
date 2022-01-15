package debug

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/oakmound/grove/components/fonthelper"
	"github.com/oakmound/grove/components/textqueue"

	"github.com/oakmound/game-template/internal/layers"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
	"golang.org/x/image/colornames"
)

var mainIsRunning bool

var Scene = scene.Scene{
	Start: func(ctx *scene.Context) {
		mainIsRunning = true

		dFPS := render.NewDrawFPS(0.1, nil, 900, 10)
		ctx.DrawStack.Draw(dFPS, 1)
		fnt := render.DefaultFont()

		titleFnt, _ := fnt.RegenerateWith(fonthelper.WithSizeAndColor(20, image.NewUniform(colornames.Crimson)))

		// Arbitrary locations for where we will draw our help text in the debug scene
		keyX := 10.0
		tqX := 140.0

		kTitle := titleFnt.NewText("KeyUsage", keyX, 10)
		ctx.DrawStack.Draw(kTitle, 0)

		// vertical line to seperate the 2 sections.
		bound2 := render.NewLine(tqX-10, 5, tqX-10, float64(ctx.Window.Height()), colornames.Beige)
		ctx.DrawStack.Draw(bound2, 0)

		logTitle := titleFnt.NewText("Logging", tqX, 10)
		ctx.DrawStack.Draw(logTitle, 0)

		fntSize := 14.0
		fnt, _ = fnt.RegenerateWith(fonthelper.WithSizeAndColor(fntSize, image.NewUniform(colornames.Aliceblue)))

		bkg := render.NewColorBox(50, 50, colornames.Lightslategray)
		ctx.Window.(*oak.Window).SetBackground(
			bkg.Modify(mod.Resize(ctx.Window.Width(), ctx.Window.Height(), mod.CubicResampling)),
		)

		tq := textqueue.New(
			ctx, "",
			floatgeom.Point2{tqX, 50}, layers.Hover,
			fnt.Copy(),
			15*time.Second,
		)
		ctx.DrawStack.Draw(tq, 0)

		keyQFnt, _ := fnt.RegenerateWith(fonthelper.WithColor(image.NewUniform(colornames.Aliceblue)))

		// This will be hooked into our local down but then we will
		keyQueue := textqueue.New(
			ctx, "KeyChange",
			floatgeom.Point2{keyX, 50}, layers.Hover,
			keyQFnt,
			4*time.Second,
		)
		ctx.DrawStack.Draw(keyQueue, 0)

		dbgStart := time.Now()
		go func() {
			d := 5.0 * time.Second
			t := time.NewTimer(d)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					ctx.EventHandler.Trigger("DisplayText", fmt.Sprintf("<--- %.0f seconds after starts", time.Now().Sub(dbgStart).Seconds()))
					t.Reset(d)
				case <-ctx.Done():
					mainIsRunning = false
					return
				}
			}
		}()

		// start the lossy event watcher
		startWatchdog(ctx)

		// Recreate main event handlers here.
		ctx.EventHandler.GlobalBind(RecreationNeeded, func(id event.CID, payload interface{}) int {

			event.GlobalBind(key.Down, func(_ event.CID, k interface{}) int {
				kValue := k.(key.Event)
				ctx.EventHandler.Trigger("KeyChange", fmt.Sprintf("Press: %v", keyCodeString(kValue.Code)))
				return 0
			})
			event.GlobalBind(key.Up, func(_ event.CID, k interface{}) int {
				kValue := k.(key.Event)
				ctx.EventHandler.Trigger("KeyChange", fmt.Sprintf("Release: %v", keyCodeString(kValue.Code)))
				return 0
			})

			return 0
		})

		return

	}, Loop: func() (cont bool) {
		// Keep running if main loop is running
		return mainIsRunning
	}, End: func() (nextScene string, result *scene.Result) {
		// There is never anything after this scene so just exit
		os.Exit(1)
		return "", nil
	},
}

func keyCodeString(c key.Code) string {
	s := c.String()
	return s[4:]
}

type displayWriter struct {
	ctx *scene.Context
}

func (dw *displayWriter) Write(data []byte) (n int, err error) {
	dw.ctx.EventHandler.Trigger("DisplayText", string(data))
	os.Stdout.Write(data)
	return 1, nil
}

func standardGoToBind(ctx *scene.Context, cmdKey, sceneName string) btn.Option {
	return btn.And(
		btn.Click(event.Empty(func() { go ctx.Window.GoToScene(sceneName) })),
		btn.Binding(key.Down+cmdKey, event.Empty(func() { go ctx.Window.GoToScene(sceneName) })),
	)
}
