<p align="center">
  <img src="https://github.com/kilianc/pretender/assets/385716/3344aed5-e974-4402-806b-c1386d201469" height="150">
</p>

<p align="center">
  <img src="https://github.com/kilianc/pretender/actions/workflows/go.yml/badge.svg?branch=main">
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
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ v1.1.0

• starting server on port 8080
• using responses file: examples/example.json
• press ctrl+c to stop
````

### Install

```sh
go install github.com/kilianc/pretender/cmd/pretender@v1.1.0
```

### Usage

Every line in `examples/example.json` will match one consecutive http response when hitting `http://localhost:8080`

```sh
pretender --port 8080 --responses examples/example.json
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
    "delay_ms": 1000
  },
  {
    "body": "{\"hello\":\"world\"}",
    "headers": {"Content-Type":"application/json"},
  },
  // ...
]
```

#### A valid response definition can contain the following fields

| name          | description                            | default                         |
| ------------- | -------------------------------------- | ------------------------------- |
| `status_code` | HTTP status code                       | `200`                           |
| `body`        | HTTP response body                     | `""`                            |
| `headers`     | HTTP headers                           | `{"content-type":"text/plain"}` |
| `delay_ms`    | Number of ms to wait before responding | `0`                             |

### Local Development

These are the usual suspects

```sh
make run
make build
make test
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

MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/main/LICENSE)
