# Environment Variables

## Templating functions

* `requiredEnv`: allows you to declare a particular environment variable as required for template rendering. If the value is not set, the template rendering step will fail with an error message.
* `env`: pulls the value out an existing environment variable. Returns `""` if the environment variable is not set

### Note

If you wish to treat your environment variables as strings always, even if they are boolean or numeric values you can use `{{ env "ENV_NAME" | quote }}` or `"{{ env "ENV_NAME" }}"`. These approaches also work with the `requiredEnv` function.

## .env

appfile supports loading environment variables from a `.env` file. By default, it looks for a `.env` file in the current working directory. The location can be customized with the `--env-file` option.