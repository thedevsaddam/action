# Reviewpad Action Release Guide

This document describes how we version and release the GitHub Action.

The goal of our release is to have stable codenamed versions, denoted by `vNUM.x`, in such a way that the users of the action do not have to constantly bump the version of their workflow to benefit of new features.

This repository contains two artifacts:

1. GitHub action. The GitHub action executes the Docker image [reviewpad/action](https://hub.docker.com/repository/docker/reviewpad/action) where the arguments are passed as environment variables.
   
2. Golang executable and Docker distribution. This executable imports the open source [main Reviewpad repository](https://github.com/reviewpad/reviewpad) and also repositories which are currently private. 

On each push, a workflow is triggered to automatically build and push a new version of the Docker image. This new version is then used in the [canary branch](https://github.com/reviewpad/action/tree/canary) so that we can test experimental versions of the action:

```yaml
name: Reviewpad Action

on: pull_request

jobs:
  calls-reviewpad:
    runs-on: ubuntu-latest
    steps:
      - name: Running reviewpad
        uses: reviewpad/action@canary
```

We use [semantic versioning](https://semver.org/) for the main Go executable. Its release occurs in the commits with the message `Update reviewpad to vX.Y.Z` (e.g. [v1.0.0](https://github.com/reviewpad/action/commit/bb1d889ac9ef53627ff0eaae48ee242994b32811)).

The process of that release is manually done by an official Reviewpad contributor and amounts to tag the commit image with the semantic version and push it.

Each release also updates the stable version through either a push of a new stable version or a re-tag and push of the current stable version. 

It should be transparent through the digest of the Docker images to identify for any given stable version, the actual version and which commit in this repository originated the Docker image.

![docker-digests](https://user-images.githubusercontent.com/601882/174272063-76c22e36-1a32-4de2-b826-c3f9a50ace54.png)


Starting with `v1.x`, we only release the stable versions of the GitHub action. 
These versions point to an evolving Docker image with the same version. 
This is the way we can update the version of Reviewpad you are using without requiring manual version bumps of the action.
