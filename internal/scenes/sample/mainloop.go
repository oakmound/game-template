package sample

import (
	"image/color"

	"github.com/oakmound/game-template/internal/scenes"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

// Scene that our game starts on.
var Scene = scene.Scene{
	Start: func(ctx *scene.Context) {
		rb := render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255})
		center := floatgeom.Point2{float64(ctx.Window.Bounds().X()),
			float64(ctx.Window.Bounds().Y())}.DivConst(2)
		rb.SetPos(center.X(), center.Y())
		ctx.DrawStack.Draw(rb, 1, 0)
	},
	End: func() (nextScene string, result *scene.Result) {
		return scenes.Sample, nil
	},
}
