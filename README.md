# Oak Game Template

This repo is meant to provide a basic starting point for an Oak game.
You should feel more than free to customize everything.

In addition to providing a scaffolding there is a set of mage (the rake like go build tool) scripts.
These scripts are meant to help get you up to speed even faster but do require mage to run.

## Install Mage

Get mage from [here](https://magefile.org/), or:

``` bash
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

Ensure $GOPATH/bin is in your PATH if, following this, the mage executable cannot be found.

## Run the Mage Targets under bootstrap

To replace your project name run
`mage bootstrap:replaceProjectName <username>/<repository_name>`

For example, `mage bootstrap:replaceProjectName oakmound/game-template`

## Where to Start

1) See internal\scenes\sample\mainloop.go and start making your first scene!
2) Check out the places with "CONSIDER:" and see if you want to perform the follow up steps. These locations are generally a bit out of the way (ie not in the sample scene).
3) Setup a project in Itch and connect your github actions pipeline to push builds to it. Follow the steps in .github\workflows\README.md to get your itchupload pipe running!
4) Remove or adjust the license to your project's needs.
