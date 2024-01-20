cmd/registerPlugins.go:
	./cmd/build_registerPlugins.bash > cmd/registerPlugins.go

clean:
	rm cmd/registerPlugins.go

goabe:
	go generate .
	go build .
