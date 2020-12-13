# appfile ![CI](https://github.com/renehernandez/appfile/workflows/CI/badge.svg)

Deploy App Specs to DigitalOcean App Platform

## About

`appfile` is a declarative spec for deploying apps to the DigitalOcean App Platform. It lets you:

* Keep a directory of app spec values files and maintain changes in version control
* Apply CI/CD to configuration changes
* Visualize a diff of the changes to be applied

## Installation

`appfile` can be installed in several ways.

### Homebrew

You can install directly using the `renehernandez/taps` as follows:

```shell
$ brew install renehernandez/taps/appfile
```

### Download releases

You can always download the released binaries directly from the Github Releases page. For the latest releases check [here](https://github.com/renehernandez/appfile/releases)

### Github Action

You can leverage `appfile` with your Github Actions workflows, by using `action-appfile`:

* Marketplace: https://github.com/marketplace/actions/github-action-for-appfile-cli
* Repository URL: https://github.com/renehernandez/action-appfile

## Getting Started

Let's look at an example and see how `appfile` can help you to manage your App specification and deployments.

## Deploy a web service

This example deploys an App containing a service definition. The 2 environments: *review* and *production* will customize the final specification of the app to be deployed. Let's look at the `appfile.yaml`, `app.yaml` and environments definitions below.

```yaml
# appfile.yaml
environments:
  review:
  - ./envs/review.yaml
  production:
  - ./envs/production.yaml

specs:
- ./app.yaml
```

```yaml
# app.yaml
name: {{ .Values.name }}

services:
- name: web
  github:
    repo: <repo-url>
    branch: main
    deploy_on_push: {{ .Values.deploy_on_push }}
  envs:
  - key: WEBSITE_NAME
    value: {{ requiredEnv "WEBSITE_NAME" }}
```

```yaml
# review.yaml
name: sample-review

deploy_on_push: true
```

```yaml
# production.yaml
name: sample-production

deploy_on_push: false
```

You can deploy your App in review by running:

```console
WEBSITE_NAME='Appfile Review' appfile sync --file /path/to/appfile.yaml --environment review
```

The final App spec to be synced to DigitalOcean would be:

```yaml
name: sample-review

services:
- name: web
  github:
    repo: <repo-url>
    branch: main
    deploy_on_push: true
  routes:
  - path: /
  envs:
  - key: WEBSITE_NAME
    value: Appfile Review
```

Or you can deploy your App in production:

```console
WEBSITE_NAME='Appfile Prod' appfile sync --file /path/to/appfile.yaml --environment production
```

The final App spec to be synced to DigitalOcean would be:

```yaml
name: sample-production

services:
- name: web
  github:
    repo: <repo-url>
    branch: main
    deploy_on_push: false
  routes:
  - path: /
  envs:
  - key: WEBSITE_NAME
    value: Appfile Prod
```

To learn more about `appfile`, check out the [docs](https://renehernandez.github.io/appfile/latest)

## Contributing

Check out the [Contributing](docs/CONTRIBUTING.md) page.

## Changelog

For inspecting the changes and tag releases, check the [Changelog](CHANGELOG.md) page

## Appreciation

This project is inspired in [helmfile](https://github.com/roboll/helmfile), from which I have borrowed heavily for the first iteration.

## License

Check out the [LICENSE](LICENSE) for details.
