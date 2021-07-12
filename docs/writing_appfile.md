# Writing appfile

## Defaults

* The default name for an appfile is `appfile.yaml`
* The default environment is `default`
* The access token to DigitalOcean can be specified through the `access-token` option or the `DIGITALOCEAN_ACCESS_TOKEN` environment variable

## Templating

Appfile uses [go templates](https://godoc.org/text/template) for templating your `appfile.yaml`. While golang ships several built-in functions, we have added all of the functions in the [sprig library](https://godoc.org/github.com/Masterminds/sprig).

We also added the following functions:

* `requiredEnv`: allows you to declare a particular environment variable as required for template rendering. If the value is not set, the template rendering step will fail with an error message.
* `toYaml`: allows you to get a values block and output the corresponding yaml representation

## Environment Variables

Environments variables can be used anywhere for templating the appfile.

## Paths Overview

Using spec files in conjunction with CLI arguments can be a bit confusing.

A few rules to clear up this ambiguity:

* Absolute paths are always resolved as absolute paths
* Relative paths referenced in the appfile spec itself are relative to that spec.
* Relative paths referenced on the command line are relative to the current working directory the user is in
