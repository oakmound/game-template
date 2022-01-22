# Github Workflows

This folder contains sample workflows which may take some bootstrapping on your part.
The idea is that these can be run on certain events such as pushes, merges, or just run from webconsole.

[Actions guide](https://docs.github.com/en/actions/learn-github-actions)

## itchupload

This workflow is meant to allow for an easy pipeline from your project to a given itch project.
By default it will build and push to itch whenever you push or create a pull to the main branch.
If not using modules you will need to uncomment out the GOPATH info and change any references from GOROOT to GOPATH.
Derived as much of this project was from the oakmound project that can be found [here](https://oakmound.itch.io/dashking).

## Setup itch flow

- Get an Itch Account
- Create an [itch project](https://itch.io/game/new) This will get you a project url of the form `https://<username>.itch.io/<gamename>`
- Update the itchupload.yaml and replace youritchusername/youritchproject with `<username>/<gamename>`
- If not a private github repo consider locking it down in one of the following ways to protect your secret. While it is masked by default a contributor could push a different version of the action if they wanted
  - Lock down who can push directly to your main branch (security against exfiltration)
  - Hide your actions tab (`https://github.com/<username>/<repository_name>/settings/actions`)
  - Require Approval for outside contributors in Actions (`https://github.com/<username>/<repository_name>/settings/actions`)
- If you have not previously uploaded your itch credentials to your repo or org as the secret BUTLER_API_KEY do the following
  - [Generate an API key](https://itch.io/user/settings/api-key) and [retreive its value](https://itch.io/docs/butler/login.html)
  - Go to `https://github.com/<username>/<repository_name>/settings/secrets/actions`
  - Add a new repository secret with the name BUTLER_API_KEY with the value you got from the api key page
- For the in browser component, remember to link the uploaded build from github actions to be embedded within itch.
  - Keep in mind too that the embedded canvas will need some buffer space, something around 10px w and h.
  