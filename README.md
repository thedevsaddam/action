# Reviewpad GitHub Action

This action runs the docker image [exploredev/action](https://hub.docker.com/repository/docker/exploredev/action).

It reads and automates the pull request workflows specified in the `revy.yml` file at the root of the GitHub repository.

For more information, check out the [official documentation](https://docs.reviewpad.com).

## Inputs

- **repository**: Uses default `${{ github.repository }}`
- **prnumber**: Uses default `${{ github.event.pull_request.number }}`
- **token**: Uses default `${{ github.token }}`

## Outputs

None.

## Usage examples

**This action should only be used on `pull_request` related events.**

Add the following step to a GitHub Action job:

```yaml
- name: Run reviewpad action
  uses: reviewpad/action@v0.0.4
```


By default this action uses the `github-actions[bot]` PAT.

As described in the [official GitHub documentation](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow):

> When you use the repository's GITHUB_TOKEN to perform tasks, events triggered by the GITHUB_TOKEN will not create a new workflow run.

If you want to use more advanced features such as the auto-merge feature, we recommend that you explicitly pass a PAT to run this action:

```yaml
- name: Run reviewpad action
  uses: reviewpad/action@v0.0.4
  with:
    token: ${{ secrets.GH_TOKEN }}
```

[Please follow this link to know more](https://docs.reviewpad.com/docs/install-github-action-with-github-token).