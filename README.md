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

You can install directly with using the `renehernandez/taps` as follows:

```shell
$ brew install renehernandez/taps/appfile
```

### Download releases

You can always download the released binaries directly from the Github Releases page. For the latest releases check [here](https://github.com/renehernandez/appfile/releases)

### Github Action

You can leverage `appfile` with your Github Actions workflows, by using `action-appfile`:

* Marketplace: https://github.com/marketplace/actions/github-action-for-appfile-cli
* Repository URL: https://github.com/renehernandez/action-appfile

## Defaults

* The default name for an appfile is `appfile.yaml`
* The default environment is `default`
* The access token to DigitalOcean can be specified through the `access-token` flag or the `DIGITALOCEAN_ACCESS_TOKEN` environment variable

## Getting Started

Let's look at several app examples and see how `appfile` can help you to manage your App specification and deployments.

### Introductory example

This example will deploy a static site without any custom values nor environment values.

```yaml
# appfile.yaml
specs:
- ./app.yaml
```

```yaml
# app.yaml
name: sample-html

static_sites:
- environment_slug: html
  github:
    branch: main
    deploy_on_push: true
    repo: renehernandez/sample-html
  name: sample-html
```

Sync your App specification to DigitalOcean App Platform by running:

Using access token flag:

```console
appfile sync --file /path/to/appfile.yaml --access-token <token>
```

Using `DIGITALOCEAN_ACCESS_TOKEN` environment variable

```console
appfile sync --file /path/to/appfile.yaml
```

For the example above, you don't need `appfile`, you can use instead the [doctl cli](https://github.com/digitalocean/doctl) to deploy your app. Let's look at a more interesting example next with a fictitious Django app, which will show the flexibility of environments.

### Intermediate example

This next example deploys an App containing a service definition. The 2 environments: *review* and *production* will customize the final specification of the app to be deployed. Let's look at the `appfile.yaml`, `app.yaml` and environments definitions below.

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

### A more complex example

Finally let's go over a more complex scenario, using a Rails app as an example. The app spec declares a rails service, a migration job and a database. The 2 environments: *review* and *production* will customize the final App spec that gets synced with DigitalOcean. Let's look at the `appfile.yaml`, `app.yaml` and environments definitions below.


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
name: {{ .Values.name }}

services:
- name: rails-app
  image:
    registry_type: DOCR
    repository: <repo_name>
    tag: {{ requiredEnv "IMAGE_TAG" }}
  instance_size_slug: {{ .Values.rails.instance_slug }}
  instance_count: {{ .Values.rails.instance_count }}
  envs:
{{- range $key, $value := .Values.rails.envs }}
  - key: {{ $key }}
    value: {{ $value }}
{{- end }}

{{- if eq .Environment.Name "review" }}
- name: postgres
  image:
    registry_type: DOCR
    repository: postgres
    tag: '12.4'
  internal_ports:
    - 5432
  envs:
{{- range $key, $value := .Values.postgres.envs }}
  - key: {{ $key }}
    value: {{ $value }}
{{- end }}
{{- end }}

jobs:
- name: migrations
  image:
    registry_type: DOCR
    repository: <repo_name>
    tag: {{ requiredEnv "IMAGE_TAG" }}****
  envs:
{{- range $key, $value := .Values.migrations.envs }}
  - key: {{ $key }}
    value: {{ $value }}
{{- end }}

{{- if eq .Environment.Name "production" }}
databases:
- name: db
  production: true
  cluster_name: mydatabase
  engine: PG
  version: "12"
{{- end }}
```

```yaml
# review.yaml
name: sample-{{ requiredEnv "REVIEW_HOSTNAME" }}

.common_envs: &common_envs
  DB_USERNAME: postgres
  DB_PASSWORD: password
  RAILS_ENV: production

rails:
  instance_slug: basic-xxs
  instance_count: 1
  envs:
  <<: *common_envs

postgres:
  envs:
    POSTGRES_USER: postgres
    POSTGRES_DB: mydatabase
    POSTGRES_PASSWORD: password

migrations:
  envs:
  <<: *common_envs
```

```yaml
# production.yaml
name: sample-production

.common_envs: &common_envs
  DB_USERNAME: postgres
  DB_PASSWORD: strong_password
  RAILS_ENV: production

rails:
  instance_slug: professional-xs
  instance_count: 3
  envs:
  <<: *common_envs

migrations:
  envs:
  <<: *common_envs
```

You can deploy your App in review by running:

```console
IMAGE_TAG='fad7869fdaldabh23' REVIEW_HOSTNAME='fix-bug' appfile sync --file /path/to/appfile.yaml --environment review
```

This would deploy a public rails service, and internal postgres service (the database running on a container) and would run the migration job. The final App spec to be synced to DigitalOcean would be:

```yaml
name: sample-fix-bug

services:
- name: rails-app
  image:
    registry_type: DOCR
    repository: <app-repo>
    tag: fad7869fdaldabh23
  instance_size_slug: basic-xxs
  instance_count: 1
  routes:
  - path: /
  envs:
  - key: DB_PASSWORD
    value: password
  - key: DB_USERNAME
    value: postgres
  - key: RAILS_ENV
    value: production

- name: postgres
  image:
    registry_type: DOCR
    repository: postgres
    tag: '12.4'
  internal_ports:
    - 5432
  envs:
  - key: POSTGRES_DB
    value: mydatabase
  - key: POSTGRES_PASSWORD
    value: password
  - key: POSTGRES_USER
    value: postgres

jobs:
- name: migrations
  image:
    registry_type: DOCR
    repository: <migration-repo>
    tag: fad7869fdaldabh23
  envs:
  - key: DB_PASSWORD
    value: password
  - key: DB_USERNAME
    value: postgres
  - key: RAILS_ENV
    value: production
```

Or you can deploy your App in production instead:

```console
IMAGE_TAG='fad7869fdaldabh23' appfile sync --file /path/to/appfile.yaml --environment production
```

This would deploy a public rails service and a migration job. Both components would connect to an existing database. The final App spec to be synced to DigitalOcean would be:

```yaml
name: sample-production

services:
- name: rails-app
  image:
    registry_type: DOCR
    repository: <app-repo>
    tag: fad7869fdaldabh23
  instance_size_slug: professional-xs
  instance_count: 3
  routes:
  - path: /
  envs:
  - key: DB_PASSWORD
    value: strong_password
  - key: DB_USERNAME
    value: postgres
  - key: RAILS_ENV
    value: production

jobs:
- name: migrations
  image:
    registry_type: DOCR
    repository: <migration-repo>
    tag: fad7869fdaldabh23
  envs:
  - key: DB_PASSWORD
    value: strong_password
  - key: DB_USERNAME
    value: postgres
  - key: RAILS_ENV
    value: production

databases:
- name: db
  production: true
  cluster_name: mydb
  engine: PG
  version: "12"
```

You can check out more examples in the examples folder of this repo

## Writing appfile

For patterns, resources and tips writing appfile, check the [Writing appfile guide](docs/writing_appfile.md).

## CLI Reference

See [CLI Reference Documentation](docs/cli_reference.md) for information about each available command.

## Contributing

Check out the [Contributing](docs/CONTRIBUTING.md) page.

## Changelog

For inspecting the changes and tag releases, check the [Changelog](CHANGELOG.md) page

## Appreciation

This project is inspired in [helmfile](https://github.com/roboll/helmfile), from which I have borrowed heavily for the first iteration

## License

Checkout the [LICENSE](LICENSE) for details