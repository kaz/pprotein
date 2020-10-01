.PHONY: run
run:
	go run .

.PHONY: build
build: pprotein

pprotein: rice-box.go
	go build -o $@ -ldflags "-w -s"

rice-box.go:
	go get github.com/GeertJohan/go.rice/rice
	rice embed-go

.PHONY: clean
clean:
	rm -rf pprotein rice-box.go
