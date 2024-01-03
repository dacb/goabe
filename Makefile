plugins: plugins/example.so

all: $(plugins)

clean:
	rm -rf $(plugins)

%.so: %.go
	go build -buildmode=plugin -o $@ $<

