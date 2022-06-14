//go:build !prod
// +build !prod

package loading

import (
	"time"

	"github.com/oakmound/oak/v4/scene"
)

func waitInProduction(ctx *scene.Context, start time.Time) {}
