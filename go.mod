module github.com/oakmound/game-template

go 1.18

require (
	github.com/disintegration/gift v1.2.1
	github.com/magefile/mage v1.12.1
	github.com/oakmound/grove/components/fonthelper v0.0.0-20220111021726-41ed3856c8b7
	github.com/oakmound/grove/components/textqueue v0.0.0-20220111021726-41ed3856c8b7
	github.com/oakmound/oak/v4 v4.0.0-alpha.2
	golang.org/x/image v0.0.0-20220321031419-a8550c1d254a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20201218220906-28db891af037 // indirect
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/BurntSushi/xgbutil v0.0.0-20190907113008-ad855c713046 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20220320163800-277f93cfa958 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jfreymuth/pulse v0.1.0 // indirect
	github.com/oakmound/alsa v0.0.2 // indirect
	github.com/oakmound/libudev v0.2.1 // indirect
	github.com/oakmound/w32 v2.1.0+incompatible // indirect
	github.com/oov/directsound-go v0.0.0-20141101201356-e53e59c700bf // indirect
	golang.org/x/exp v0.0.0-20220414153411-bcd21879b8fd // indirect
	golang.org/x/exp/shiny v0.0.0-20220414153411-bcd21879b8fd // indirect
	golang.org/x/mobile v0.0.0-20220325161704-447654d348e3 // indirect
	golang.org/x/sys v0.0.0-20220403205710-6acee93ad0eb // indirect
)

replace (
	github.com/oakmound/grove => ..\grove
	github.com/oakmound/grove/components/fonthelper => ../grove/components/fonthelper
	github.com/oakmound/grove/components/textqueue => ../grove/components/textqueue
	github.com/oakmound/oak/v4 => ..\oak
)
