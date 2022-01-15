// Package layers provides a set of utilities for sharing a general schema of layering that should help
// provide unity between games and helper packages. These are recommendations that have helped in the past.
package layers

// Layer levels defined for general reuse and resizing
const (
	Back        = 0
	Front       = 100
	Hover       = 200
	ModalBottom = 300
	ModalMid    = 400
	ModalTop    = 50
	Debug       = 600
)

// Naming for the draw stack layers. See the actual implementation in the top level main.go
const (
	StackBackground = 0
	StackDynamic    = 1
	StackUI         = 2
)
