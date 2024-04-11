# `pretender`

A naive HTTP mock server that responds to requests with a sequence from a text file. Used for controlling e2e testing.

```
❯ bin/pretender

██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗
██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗
██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ v1.0.0

• starting server on port 8080
• using responses file: responses.txt
• press ctrl+c to stop
````

### Usage

To use default value `8080` for `port` and `README.md` for `responses`, run

```
make run
```

To build the binary and test different arguments, run

```
make build
bin/pretender --port 8080 --responses responses.txt
```

### Testing

```
make test
```
