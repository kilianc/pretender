# Release Process

Pick a new version:

```sh
export NEW_VERSION="v1.x.x"
```

Update all version strings:

```sh
make version-bump version=${NEW_VERSION}
```

Update `CHANGELOG.md`:

```sh
make changelog tag=${NEW_VERSION}
```

Check changes and commit/push:

```sh
git status
git add .
git commit -m "chore: bump version to ${NEW_VERSION}"
```

You should see something similar to:

```console
On branch feat/foobar
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   CHANGELOG.md
	modified:   Makefile
	modified:   README.md
	modified:   cmd/pretender/main.go
```

Create a new GitHub release using the value in `${NEW_VERSION}` as name and tag. The new release will trigger the automation to upload the binaries and push a docker image to DockerHub.
