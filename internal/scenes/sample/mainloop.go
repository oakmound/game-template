package sample

import (
	"image/color"

	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var stayInScene bool

var Scene = scene.Scene{
	Start: func(ctx *scene.Context) {
		rb := render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255})
		rb.SetPos(float64(ctx.Window.Width())/2, float64(ctx.Window.Height())/2)
		ctx.DrawStack.Draw(rb, 1, 0)

	},
	Loop: scene.BooleanLoop(&stayInScene),
	End: func() (nextScene string, result *scene.Result) {
		stayInScene = true
		return scenes.Sample, nil
	},
}
