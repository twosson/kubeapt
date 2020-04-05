# Kubeapt

Kubernetes dashboard for developers

## Running 

`$ apt dash`

## Developing

### Prerequisites

* Go 1.14
* npm
* yarn

### Environment variables

* `DASH_DISABLE_OPEN_BROWSER`  - set to a non-empty value if you don't the browser launched when the dashboard start up.
* `DASH_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)

### Running development web UI

`$ make setup-web`

### Building binary with embedded web assets

Create `./build/apt`: `$ make apt-dev`
