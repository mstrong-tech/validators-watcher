# Validators Watcher

## Requierements

* Golang v1.20
* Access to a Beacon API for the networks to be monitored

## Instructions

1. Build watcher cli with `go build -o watcher cmd/main.go`.
2. Write your watcher configuration based on the [example config](config.yaml.example) file.
3. Run your watcher binary pointing to your modified configuration file with `./watcher --config path/to/config.yaml`.

## CLI Options

```
Usage of ./watcher:
  -config string
        path to config file (default "config.yaml")
  -log-level uint
        log level (default 4)
```
