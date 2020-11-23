# CLI Reference

```shell
Deploy app platform specifications to DigitalOcean

Usage:
  appfile [command]

Available Commands:
  destroy     Destroy apps running in DigitalOcean
  diff        Diff local app spec against app spec running in DigitalOcean
  help        Help about any command
  sync        Sync all resources from app platform specs to DigitalOcean

Flags:
  -t, --access-token string   API V2 access token
  -e, --environment string    root all resources from spec file (default "default")
  -f, --file string           load appfile spec from file (default "appfile.yaml")
  -h, --help                  help for appfile
      --log-level string      Set log level (default "info")
  -v, --version               version for appfile

Use "appfile [command] --help" for more information about a command.
```