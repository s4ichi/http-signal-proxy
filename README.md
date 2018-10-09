# http-signal

http-signal is the tiny process to expose a process that traps unix signals to http.

## Installation
See https://github.com/everysick/http-signal/releases to install binary.

## Usage

You can use following options.

```console
$ ./http-signal --help
Usage:
  http-signal [OPTIONS]

Application Options:
  -c, --command= Destination command to proxying signals
  -p, --port=    Port number to listen http
  -f, --prefix=  Prefix path of url to listen http (default: /http-signal)

Help Options:
  -h, --help     Show this help message
```

### Signals
Now, http-signal supports `INT, TERM, HUP, QUIT, WINTH, USR1, USR2, TTIN, TTOU`.

### Endpoints
You can request with http (e.g., `http://<YOUR_HOSTNAME>:<PORT>/<PREFIX>/sigint`)

## Setup

```console
$ make setup
```

## License
MIT

## Author
Takamasa Saichi
