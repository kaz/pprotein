.PHONY: noop
noop:

.PHONY: run
run:
	go run ./cli/pprotein

.PHONY: run-agent
run-agent:
	go run ./cli/pprotein-agent

.PHONY: build
build: pprotein pprotein-agent

pprotein: view/dist
	go build -ldflags="-w -s" -gcflags="-trimpath=$$PWD" -asmflags="-trimpath=$$PWD" ./cli/pprotein

pprotein-agent:
	go build -ldflags="-w -s" -gcflags="-trimpath=$$PWD" -asmflags="-trimpath=$$PWD" ./cli/pprotein-agent

view/dist:
	npm --prefix view ci
	npm --prefix view run build

.PHONY: clean
clean:
	rm -rf pprotein pprotein-agent view/dist
