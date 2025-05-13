# euvd-go

A Go SDK and commandline application to pull up data from the [ENISA EU Vulnerability Database](https://euvd.enisa.europa.eu/).

## Running

Run from binary:
```shell
% euvd search -vendor=f5
...
% euvd lookup EUVD-2025-14349
...
```

## Building

```shell
% make build
```

## Usage as an SDK

Take a look at our own `cli.go` [here](https://github.com/hazcod/euvd-go/blob/main/cmd/cli.go).
