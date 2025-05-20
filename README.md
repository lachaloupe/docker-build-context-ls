# docker-context-ls

> Make & Docker, a love story with fewer misunderstandings.

# Context

This is a simple standalone tool that lists all files that will be included in a Docker build context.
It can help avoid the build context creation overhead when there are no file changes.

# Install

If golang is installed:

```bash
$ go install github.com/lachaloupe/docker-context-ls@latest
```

# Usage

Run with the context directory as the only argument.
The tool will look for a `.dockerignore` file inside this directory.

```bash
$ docker-context-ls .
```

# Makefile

The motivation for this tools is to use it inside a `Makefile` such as this one:

```Makefile
IMAGE=$(USER)/app

.image.done: $(shell docker-context-ls .)
	docker build --tag $(IMAGE) .
    touch $@

all: .image.done
    docker run --rm $(IMAGE)
```

Here, the `.image.done` target depends on all files that will be in the context.
Thus, the docker context will be created using `make` usual rules i.e. only if one of the dependencies is out of date.
Of course, when that happens, the usual docker build cache will rebuild only the necessary layers.
But if creating the build context is significant, this can be helpful in making operations more efficient.

# Using `.dockerignore`

One approach to create minimal contexts is to exclude everything and be intentional about what is included.

```
# example for typical python projects
*
!pyproject.toml
!src
!tests
**/__pycache__
**/*.egg-info
**/*.dist-info
```
