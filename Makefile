all: goabe

clean:
	rm goabe
	cd cmd; make clean

goabe:
	go generate ./cmd
	go build .

test: goabe
	./goabe --threads 10 run --steps 100
