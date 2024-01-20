# goabe
go agent based engine

## Getting started
Install go version 1.21 or later. Make sure the cobra generator is installed, e.g.
```
go install github.com/spf13/cobra-cli@latest
```
To generate your starting config file, do something like:
```
go run main.go config create
```
Now you can edit `goabe.json` or rename it and reference it with the `--config` option, e.g.,
```
go run main.go --config myconfig.json run --steps 10 --threads 10
```

### Plugins
To build the plugins, first use:
```
make build-plugins
```
This will build the `.so` files in the `plugins` directory.  This directory needs to be searchable
and readable by goabe.

###
To download the LIF files, I used:
```
wget http://www.ibiblio.org/lifepatterns/lifep.zip
mkdir lifep
cd life
unzip ../lifep.zip 
```
