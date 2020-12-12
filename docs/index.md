# What is Appfile?

`appfile` is a CLI to deploy App Specs to DigitalOcean App Platform.

## Features

* Keep a directory of app spec values files and maintain changes in version control.

This allows you to support multiple environments with different components and configurations per `app.yaml` file.

* Apply CI/CD to configuration changes.

Using `appfile` directly or the [appfile action](https://github.com/renehernandez/action-appfile), you can automate the deployments of different environments based on the branch you are running. For example, review environments with lower requirements in terms of CPU and memory vs production environment.

* Visualize a diff of the changes to be applied.

It lets you print a diff in console, which helps you verify that the correct changes are going to be deployed into App Platform
