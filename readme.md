# Github Labeler

The purpose of this application is to sync your default labels with a repository hosted by Github.
Labeler automatically merges your labels into the repository using a template file.

After the run, missing labels are added, unknown labels are removed and existing labels are checked for changes.

Please note that label links to issues are lost when a label is removed.

# Usage

Just run the binary and define the following arguments:

-   `-t` Your Github Token used for Github API requests
-   `-o` Your Github Username
-   `-r` The related repository name

```sh
labeler merge -t <token> -o <owner> -r <repo>
```

Optional arguments:

-   `--dry-mode` Activates dry mode (sumulation, no changes are made)
-   `--skip-delete` Skips deletion of unknown labels

You can also use the following environment variables to define default values.

-   `LABELER_TOKEN`
-   `LABELER_OWNER`
-   `LABELER_REPO`

# Define labels

Labeler uses a JSON template file to merge the labels.
The file must be placed under `~/.config/labeler/labels.json`.

See [labels.json](labels.json) for an example.

# Development

## Build

Currently the app has not been added to any package manager.

If you want to use the app, you need to install and build it yourself.

```sh
go build ./cmd/labeler
```

After that, just copy the binary to your preferred location.

```sh
cp labeler /usr/bin/
```

## Debugging

To set arguments in debug mode, create an `.env` file in the root of the repository.

Then add the following content to it:

```
LABELER_TOKEN=<GH API bearer token>
LABELER_OWNER=<GH username>
LABELER_REPO=<GH repository name>
```

> If the file exists, VSCode applies the parameters as environment variables for the running process.
>
> See [.vscode/launch.json](.vscode/launch.json)

Alternatively, you can override the parameter by setting the args section in the `.vscode/launch.json` file.

Please make sure you never push that file to Github.
Because of this, the `.env` file is already excluded in the `.gitignore` file.
