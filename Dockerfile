# --------------------------------------------------

FROM alpine AS pprotein

RUN apk add go npm make

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN make build

# --------------------------------------------------

FROM alpine AS tools

RUN apk add go

ENV GOPATH /go
RUN go install github.com/tkuchiki/alp/cli/alp@latest
RUN go install github.com/tkuchiki/slp/cmd/slp@latest

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache bash perl perl-dbd-mysql perl-time-hires graphviz

COPY --from=pprotein /go/src/app/pprotein /usr/local/bin/
COPY --from=pprotein /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=tools /go/bin/alp /usr/local/bin/
COPY --from=tools /go/bin/slp /usr/local/bin/

RUN mkdir -p /opt/pprotein
WORKDIR /opt/pprotein

ENTRYPOINT ["pprotein"]
