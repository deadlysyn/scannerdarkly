.PHONY: build clean

build:
	go build -o d .

clean:
	rm -f d
