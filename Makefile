plugins: plugins/example/example.so

all: $(plugins)

clean:
	rm -rf $(plugins)

%.so: %.go plugins/plugins.go Makefile
	go build -buildmode=plugin -o $@ $<

