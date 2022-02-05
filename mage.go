//go:build mage
// +build mage

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Aliases for some of the more used commands.
var Aliases = map[string]interface{}{
	"dr": DebugRun,
}

// DebugRun of the main game.
func DebugRun() error {
	err := sh.Run("go", "run", "./main.go", "--debug=true")
	return err
}

// Build namespace for general build tooling.
type Build mg.Namespace

// All distros should be built.
func (Build) All() error {
	os.Chdir("build")
	return sh.RunV("go", "run", "build.go", "-win", "-js", "-nix", "-osx", "-v")
}

// Js is a convience for the js build run pattern.
func (Build) Js() error {
	os.Chdir("build")
	toRun := []string{"run", "build.go", "-js", "-win=false"}
	fmt.Println("Running: go ", toRun)
	err := sh.RunV("go", toRun...)
	os.Chdir("..")
	return err
}

// Bootstrap tooling for additional project setup after the clone of the template repo.
// CONSIDER: Running all the bootstrap functions.
type Bootstrap mg.Namespace

// ReplaceProjectName across the template repo.
func (Bootstrap) ReplaceProjectName(newName string) {
	out, _ := sh.OutCmd("git", "config", "--get", "remote.origin.url")()
	fmt.Println("Current git path is", out)

	// validate naming
	splitName := strings.Split(newName, "/")
	if len(splitName) != 2 {
		fmt.Print("You provided \"", newName, "\" as your new name which is not of the form <username>/<project>.",
			"If this was intentional edit out this check otherwise fix your input")
		return
	}
	projectName := splitName[1]
	fmt.Println("Attempting to set the project name to", newName, " ok? y/n")
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	trimText := strings.ToLower(strings.TrimSpace(text))
	if "y" != trimText {
		fmt.Println("Stopping conversion as entry was", trimText, " and not y")
		return
	}
	fmt.Println("Starting to replace template name with yours")
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			// Get the file contents
			body, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Failed to read ", path, err, "continuing")
				return nil // still want to mutate as much as possible so continue to walk file sys
			}

			newBody := strings.ReplaceAll(string(body), "oakmound/game-template", newName)
			newBody = strings.ReplaceAll(newBody, "sample-project", projectName)

			// lazy way to not touch files that we dont need to
			// slow but its just a bootstrap and at least paranoia means that we wont affect asset files and the like
			if newBody == string(body) {
				return nil
			}

			// rewrite the contents
			os.WriteFile(path, []byte(newBody), 0666)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
