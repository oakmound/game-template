package debug

import (
	"fmt"
	"image"
	"io"
	"os"
	"time"

	"github.com/oakmound/game-template/internal/layers"
	"github.com/oakmound/grove/components/fonthelper"
	"github.com/oakmound/grove/components/textqueue"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"
	"golang.org/x/image/colornames"
)

var mainIsRunning bool

// ExtraBus to watch (main bus in this case)
const ExtraBus = "mainBus"

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
		bound2 := render.NewLine(tqX-10, 5, tqX-10, float64(ctx.Window.Bounds().Y()), colornames.Beige)
		ctx.DrawStack.Draw(bound2, 0)

		logTitle := titleFnt.NewText("Logging", tqX, 10)
		ctx.DrawStack.Draw(logTitle, 0)

		fntSize := 14.0
		fnt, _ = fnt.RegenerateWith(fonthelper.WithSizeAndColor(fntSize, image.NewUniform(colornames.Aliceblue)))

		bkg := render.NewColorBox(50, 50, colornames.Lightslategray)
		ctx.Window.(*oak.Window).SetBackground(
			bkg.Modify(mod.Resize(ctx.Window.Bounds().X(), ctx.Window.Bounds().Y(), mod.CubicResampling)),
		)

		keyQFnt, _ := fnt.RegenerateWith(fonthelper.WithColor(image.NewUniform(colornames.Aliceblue)))

		keyQueue := textqueue.New(
			ctx, []event.UnsafeEventID{},
			floatgeom.Point2{keyX, 50}, layers.Hover,
			keyQFnt,
			4*time.Second,
		)

		// Extra the reference to the main event loop so we can listen on it
		extraBusI := ctx.Value(ExtraBus)
		extraBus, _ := extraBusI.(event.Handler)
		// bind to global id of 0 as our queue is not in the main event loop so this is a nice and easy redirecter.
		extraBus.PersistentBind(key.AnyDown.UnsafeEventID, 0,
			func(cidInOtherMap event.CallerID, handler event.Handler, payload interface{}) event.Response {
				press, _ := payload.(key.Event)
				return textqueue.PrintBind(keyQueue, key.AllKeys[press.Code])
			})
		ctx.DrawStack.Draw(keyQueue, 0)

		// Create a queue to detail the state of the game.
		logQ := textqueue.New(
			ctx, []event.UnsafeEventID{},
			floatgeom.Point2{tqX, 50}, layers.Hover,
			fnt.Copy(),
			15*time.Second,
		)
		ctx.DrawStack.Draw(logQ, 0)
		// Post timing info to the loggin queue
		dbgStart := time.Now()
		go func() {
			d := 5.0 * time.Second
			t := time.NewTimer(d)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					event.TriggerForCallerOn(ctx, logQ.CID(), textqueue.TextQueuePublish,
						fmt.Sprintf("<--- %.0f seconds after start", time.Now().Sub(dbgStart).Seconds()))

					t.Reset(d)
				case <-ctx.Done():
					mainIsRunning = false
					return
				}
			}
		}()

		out := os.Stdout
		mw := io.MultiWriter(out, newWriter(ctx, logQ.CallerID))
		r, w, _ := os.Pipe()
		os.Stdout = w
		os.Stderr = w

		go func() {
			// copy all reads from pipe to multiwriter, which writes to stdout and file
			_, _ = io.Copy(mw, r)
		}()

		return
	},
	End: func() (nextScene string, result *scene.Result) {
		// There is never anything after this scene so just exit
		os.Exit(1)
		return "", nil
	},
}

func keyCodeString(c key.Code) string {
	s := c.String()
	return s[4:]
}

type bufferedEventWriter struct {
	ctx       *scene.Context
	theCaller event.CallerID
}

func newWriter(ctx *scene.Context, caller event.CallerID) bufferedEventWriter {
	return bufferedEventWriter{ctx, caller}
}

func (bew bufferedEventWriter) Write(p []byte) (int, error) {
	event.TriggerForCallerOn(bew.ctx, bew.theCaller, textqueue.TextQueuePublish,
		string(p),
	)
	return len(p), nil
}
