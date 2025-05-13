# euvd-go

A Go program that exports 1Password usage, signin and audit events to Microsoft Sentinel SIEM.

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