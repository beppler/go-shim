# go-shim

Quickly add executables to your path.

This project has been inspired on [shimexe](https://github.com/lukesampson/shimexe) by [Luke Sampson](https://github.com/lukesampson)

## Build from source

To build this project from sources you need [go](https://golang.org) compiler 1.16 or newer.

Clone the repoistory and then run

```shell
go build
```

## To create shims

Follow the steps:

* Copy the `go-shim.exe` to the name of the command you want to run, for example copy it to `go.exe`.

* Create the configuration file with the same name of the executable, but with the `.ini` extension, for example `go.ini`.

* Inside the configuration file put the full path of the correct program to be executed on the key called `command` like on the example bellow.

  ```ini
  # example configuration for go-shim
  command = C:\Program Files\Go\bin\go.exe

  # optional fixed args passed to command before the args passed to shim
  args = env
  args = CC
  args = AR
  ```
