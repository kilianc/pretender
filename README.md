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
  A naive HTTP mock server with sequential responses from a text file.
  <br><br><br>
</p>

```
❯ bin/pretender

██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗
██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗
██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ v1.0.2

• starting server on port 8080
• using responses file: responses.txt
• press ctrl+c to stop
````

### Install

```
go install github.com/kilianc/pretender/cmd/pretender@v1.0.2
```

### Usage

Every line in `responses.txt` will match one http response when hitting `http://localhost:8080`

```
pretender --port 8080 --responses responses.txt
```

### Local Development

These are the usual suspects
```
make run
make build
make test
```

Binary available in the `bin/` folder

```
bin/pretender --port 8080 --responses responses.txt
```
## License

MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/main/LICENSE)
