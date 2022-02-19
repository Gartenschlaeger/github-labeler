# Github Labeler

Purpose of this application is to merge default label sets for your repository.

# Usage

Simple run the binary and define the following arguments:

-   `-t` Your Github Token used for Github API requests
-   `-o` Your Github Username
-   `-r` The related repository name

```sh
labeler -t <token> -o <owner> -r <repo>
```

You can also use the following environment variables to define defaults.

-   `LABELER_TOKEN`
-   `LABELER_OWNER`
-   `LABELER_REPO`

# Define labels

Labeler uses a json template file to merge the labels.
The file needs to be placed at `~/.config/labeler/labels.json`.

See [labels.json](labels.json) for an example.

# Development

## Build

Currently the app is not added to any package manager.

If you want to use the app you have to install go and to build it by your own.

```sh
go build ./cmd/labeler
```

Afterwards simple copy the binary to your prefered location.

```sh
cp labeler /usr/bin/
```

## Debugging

To set arguments in debug mode create a `.env` file in root of the repository.

Then add the following content to it:

```
LABELER_TOKEN=<GH API bearer token>
LABELER_OWNER=<GH username>
LABELER_REPO=<GH repository name>
```

> If this file exists VSCode will use the parameter by default from them.

Alternativelly you can also overwrite the parameter by setting args section in `.vscode/launch.json` file.

Please ensure to never commit or push it to Github.
For this reason the `.env` file is excluded in `.gitignore`.
