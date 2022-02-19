# Github Labeler

Purpose of this application is to merge default label sets for your repository.

# Usage

```sh
labeler -t <token> -o <owner> -r <repo>
```

# Development

To set parameter in debug mode create a `.env` file in root of the repository.

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
