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

```
❯ bin/pretender

██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗
██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗
██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ v1.0.3

• starting server on port 8080
• using responses file: responses.json
• press ctrl+c to stop
````

### Install

```
go install github.com/kilianc/pretender/cmd/pretender@v1.0.3
```

### Usage

Every line in `responses.json` will match one consecutive http response when hitting `http://localhost:8080`

```
pretender --port 8080 --responses responses.json
```

### Responses files

```
[
  {
    "status_code": 200,
    "body": "hello",
    "headers": {"Content-Type":"text/plain"},
    "delay_ms": 1000
  },
  {
    "body": "{\"hello\":\"world\"}",
    "headers": {"Content-Type":"application/json"},
  },
  ...
]
```

#### A valid response definition can contain the following fields

| name          | description                            | default  |
| ------------- | -------------------------------------- | -------- |
| `status_code` | HTTP Status code                       | `200`    |
| `body`        | The response body                      | `""`     |
| `headers`     | HTTP headers                           | `{}`     |
| `delay_ms`    | Number of ms to wait before responding | `0`      |

### Local Development

These are the usual suspects
```
make run
make build
make test
```

Binary available in the `bin/` folder

```
bin/pretender --port 8080 --responses responses.json
```

### Docker

If you prefer to build and run `pretender` in a docker container, just on of these commands

````
make docker-build
make docker-run
````

## License

MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/main/LICENSE)
