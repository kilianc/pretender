<p align="center">
  <img src="https://github.com/kilianc/pretender/assets/385716/3344aed5-e974-4402-806b-c1386d201469" height="150">
</p>

<p align="center">
  <img src="https://github.com/kilianc/pretender/actions/workflows/go.yaml/badge.svg?branch=main">
  <img src="https://img.shields.io/github/release/kilianc/pretender.svg">
  <img src="https://goreportcard.com/badge/github.com/kilianc/pretender">
</p>

<p>
  <h1 align="center"><code>pretender</code></h1>
</p>

<p align="center">
  A naive HTTP mock server with sequential responses from a file.
  <br><br><br>
</p>

```sh
❯ bin/pretender

██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗
██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗
██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ v1.8.0

• starting server on port 8080
• using responses file: examples/example.json
• press ctrl+c to stop
````

### Install

With `go install`

```sh
go install github.com/kilianc/pretender/cmd/pretender@v1.8.0
```

With `docker`

```sh
docker run --rm -it \
  -p 8080:8080 \
  -v $(pwd)/examples:/examples \
  kilianciuffolo/pretender:v1.8.0 --responses /examples/example.json
```

With `curl`

```sh
export CURRENT_OS=$(uname -s | tr A-Z a-z)
export CURRENT_ARCH=$(uname -m | tr A-Z a-z | sed s/x86_64/amd64/)
export TARGZ_NAME="pretender-${CURRENT_OS}-${CURRENT_ARCH}.tar.gz"
export TARGZ_URL="https://github.com/kilianc/pretender/releases/download/v1.8.0/${TARGZ_NAME}"

curl -sOL ${TARGZ_URL}
tar -xzf ${TARGZ_NAME}

echo "successfully downloaded pretender $(./pretender --version)"
```

### Usage

Every response in `examples/example.json` will match one consecutive http response when hitting `http://localhost:8080`

```sh
pretender --port 8080 --responses examples/example.json
```

The server has a default `/healthz` endpoint that responds with a `200`. If this conflicts with your mock responses, it is possible to configure it by setting the `PRETENDER_HEALTH_CHECK_PATH` environment variable.

```sh
PRETENDER_HEALTH_CHECK_PATH=/alive pretender
```

### Responses File

Both plain text and `JSON` formats are supported.

A `TEXT` file contains one response per line:

```txt
This line is the first text/plain response body with 200 status code
This line is the second text/plain response body with 200 status code
```

A `JSON` file allows more flexibility and controls:

```jsonc
[
  {
    "status_code": 200,
    "body": "hello",
    "headers": {"content-type":"text/plain"},
    "delay_ms": 1000,
    "repeat": 5
  },
  {
    "body": {
      "hello": "world"
    },
    "headers": {"Content-Type":"application/json"}
  },
  // ...
  {
    "body": "will repeat forever",
    "repeat": -1
  }
]
```

#### A valid response definition can contain the following fields

| name          | description                                          | default                         |
| ------------- | ---------------------------------------------------- | ------------------------------- |
| `status_code` | HTTP status code                                     | `200`                           |
| `body`        | HTTP response body                                   | `""`                            |
| `headers`     | HTTP headers                                         | `{"content-type":"text/plain"}` |
| `delay_ms`    | Number of ms to wait before responding               | `0`                             |
| `repeat`      | Number of times the response repeats or `-1` for `∞` | `1`                             |

### Local Development

These are the usual suspects

```sh
make run
make build
make test
make cover
```

After running `make build` the binary available in the `bin/` folder

```sh
bin/pretender --port 8080 --responses examples/example.json
```

### Docker

If you prefer to build and run `pretender` in a docker container, just on of these commands

````sh
make docker-build
make docker-run
````

## License

MIT License, see [LICENSE](https://github.com/kilianc/pretender/blob/main/LICENSE.md)
