# Examples

## Non-customized appfile

Check the code [here](https://github.com/renehernandez/appfile/tree/main/examples/static_site)

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


## Deploying a Rails application

Check the code [here](https://github.com/renehernandez/appfile/tree/main/examples/rails_app)

The app spec declares a rails service, a migration job and a database. The 2 environments: *review* and *production* will customize the final App spec that gets synced with DigitalOcean. Let's look at the `appfile.yaml`, `app.yaml` and environments definitions below.


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
