FROM alpine

RUN apk add go npm make

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN make build

FROM alpine

RUN apk add --no-cache percona-toolkit

COPY --from=0 /go/src/app/pprotein /usr/local/bin/pprotein

ENTRYPOINT ["pprotein"]
