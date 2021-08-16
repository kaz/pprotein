.PHONY: noop
noop:

.PHONY: run
run:
	go run .

.PHONY: build
build: pprotein

pprotein: view/dist
	go build -o $@ -ldflags="-w -s" -gcflags="-trimpath=$$PWD" -asmflags="-trimpath=$$PWD"

view/dist:
	npm --prefix view ci
	npm --prefix view run build

.PHONY: clean
clean:
	rm -rf pprotein rice-box.go
