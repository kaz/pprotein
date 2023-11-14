# --------------------------------------------------

FROM golang:alpine AS pprotein

RUN apk add npm make

WORKDIR $GOPATH/src/app
COPY . .

RUN make build

# --------------------------------------------------

FROM golang:alpine AS alp

RUN go install github.com/tkuchiki/alp/cmd/alp@latest

# --------------------------------------------------

FROM golang:alpine AS slp

RUN apk add gcc musl-dev
RUN go install github.com/tkuchiki/slp/cmd/slp@latest

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache graphviz

COPY --from=pprotein /go/src/app/pprotein /usr/local/bin/
COPY --from=pprotein /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=alp /go/bin/alp /usr/local/bin/
COPY --from=slp /go/bin/slp /usr/local/bin/

RUN mkdir -p /opt/pprotein
WORKDIR /opt/pprotein

ENTRYPOINT ["pprotein"]
