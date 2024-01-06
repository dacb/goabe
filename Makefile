plugins: plugins/example/example.so plugins/example/life.so

all: $(plugins)

clean:
	rm -rf $(plugins)

%.so: %.go plugins/plugins.go Makefile
	go build -buildmode=plugin -o $@ $<

