//go:build !prod
// +build !prod

package loading

import (
	"github.com/oakmound/oak/v3/scene"
)

func waitInProduction(ctx *scene.Context) {}
