# goabe
go agent based engine

## Getting started
Install go version 1.21 or later.

### Building
To build the `goabe` binary, including running `go generate`, just use `make`!
```
make
```

### Configuration
To generate your starting config file, do something like:
```
./goabe --threads 2 config create
```
Now you can edit `goabe.json` or rename it and reference it with the `--config` option, e.g.,
```
./goabe rune --config myconfig.json run --steps 10 --threads 10
```

### Plugins
Plugins are autodiscovered by the script in `cmd/build_registerPlugins.bash`.  It expects to find
a directory with a `.go` file whose basename matches the directory name.  In that file, there needs
to be functions Init, Name, Version, Description, PreRun, PostRun.  The script just checks for those
words using pattern matching, not any examination of the go parse tree.

The command to generate the plugin registration code is autorun by `go generate` via `make`.  You
can do this yourself with:
```
./cmd/build_registerPlugins.bash
```
You probably want to redirect this to the `.go` file, e.g.,
```
./cmd/build_registerPlugins.bash > cmd/registerPlugins.go
```

Or just let `go generate` do it:
```
go generate -v ./cmd
```

## Life plugin
This is a simple example of Conway's Game of Life.  It supports arbitrary
sized matrices and standard rules.  It can read and write basic RLE files.
Note that .LIF file format was removed, but could be found in `git` history.
### RLE data
Catalog is here: https://catagolue.hatsya.com/home

The `in.rle` file contains a basic pattern with a three step repeat.
