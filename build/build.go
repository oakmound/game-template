package main

// Build cross-compiles packages on set of
// OS and architecture pairs.

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	// Defaults
	osxPairs = [][2]string{
		// I think this has to be actually run on osx
		//{"darwin", "amd64"},
		//{"darwin", "arm"},
		//{"darwin", "arm64"},
	}
	linuxPairs = [][2]string{
		//{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		//{"linux", "arm64"},
	}
	winPairs = [][2]string{
		{"windows", "386"},
		{"windows", "amd64"},
	}
	jsPairs = [][2]string{
		{"js", "wasm"},
	}
	// End Defaults
	android = [][2]string{
		{"android", "arm"},
	}

	// These are grouped together because, from my (https://github.com/200sc) perspective, they
	// are less often used by themselves. If there are valid use cases
	// to split them up into their own boolean flags then this can change.
	// I admit this is mostly because I can't think of what computer would
	// use these and would also be used for a generic program.
	nonDefaultPairs = [][2]string{
		{"dragonfly", "amd64"},
		{"freebsd", "386"},
		{"freebsd", "amd64"},
		{"freebsd", "arm"},
		{"linux", "ppc64"},
		{"linux", "ppc64le"},
		{"linux", "mips"},
		{"linux", "mipsle"},
		{"linux", "mips64"},
		{"linux", "mips64le"},
		{"netbsd", "386"},
		{"netbsd", "amd64"},
		{"netbsd", "arm"},
		{"openbsd", "386"},
		{"openbsd", "amd64"},
		{"openbsd", "arm"},
		{"plan9", "386"},
		{"plan9", "amd64"},
		{"solaris", "amd64"},
	}

	osArchPairs [][2]string

	archPairFlags = map[[2]string][]string{
		{"windows", "386"}:   {"-ldflags=-H=windowsgui"},
		{"windows", "amd64"}: {"-ldflags=-H=windowsgui"},
	}

	packageName string
	outputName  string
	verbose     bool
	useosx      bool
	usewin      bool
	uselinux    bool
	usedroid    bool
	useall      bool
	usejs       bool
	help        bool
)

func init() {
	flag.BoolVar(&verbose, "v", true, "print build commands as they are run")
	flag.StringVar(&outputName, "o", "sample-project", "output executable name")
	flag.BoolVar(&useosx, "osx", false, "build darwin executables")
	flag.BoolVar(&uselinux, "nix", false, "build linux exectuables")
	flag.BoolVar(&usewin, "win", true, "build windows exectuables")
	flag.BoolVar(&usedroid, "android", false, "build android executables")
	flag.BoolVar(&usejs, "js", false, "build js executables")
	flag.BoolVar(&useall, "all", false, "build all executables")
	flag.StringVar(&packageName, "pkg", "github.com/oakmound/game-template", "package to build")
	flag.BoolVar(&help, "h", false, "prints usage")
}

func main() {
	if help {
		fmt.Println("Usage: go run build.go <flags> -pkg <package>")
		return
	}
	flag.Parse()
	if useall {
		useosx = true
		usewin = true
		usedroid = true
		usewin = true
		usejs = true
		osArchPairs = nonDefaultPairs
	}

	if useosx {
		osArchPairs = append(osArchPairs, osxPairs...)
	}
	if uselinux {
		osArchPairs = append(osArchPairs, linuxPairs...)
	}
	if usedroid {
		osArchPairs = append(osArchPairs, android...)
	}
	if usewin {
		osArchPairs = append(osArchPairs, winPairs...)
	}
	if usejs {
		osArchPairs = append(osArchPairs, jsPairs...)
	}

	for _, pair := range osArchPairs {
		os.Setenv("GOOS", pair[0])
		os.Setenv("GOARCH", pair[1])
		buildName := outputName + "_" + pair[0] + pair[1]
		if pair[0] == "windows" {
			buildName += ".exe"
		}
		if pair[1] == "wasm" {
			buildName += ".wasm"
		}
		var out bytes.Buffer
		toRun := []string{"build", "--tags", "prod", "-o", buildName}
		if flags, ok := archPairFlags[pair]; ok {
			toRun = append(toRun, flags...)
		}
		toRun = append(toRun, packageName)
		if verbose {
			fmt.Println("Running: go ", toRun)
		}
		cmd := exec.Command("go", toRun...)
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		if verbose && out.Len() != 0 {
			fmt.Printf("%s\n", out.String())
		}
	}
}
