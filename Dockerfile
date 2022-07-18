# --------------------------------------------------

FROM alpine AS pprotein

RUN apk add go npm make

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN make build

# --------------------------------------------------

FROM alpine AS alp

RUN apk add go

ENV GOPATH /go
RUN go install github.com/tkuchiki/alp/cli/alp@latest

# --------------------------------------------------

FROM alpine AS percona-toolkit

RUN wget https://downloads.percona.com/downloads/percona-toolkit/3.4.0/binary/tarball/percona-toolkit-3.4.0_x86_64.tar.gz
RUN tar zxvf percona-toolkit-3.4.0_x86_64.tar.gz

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache bash perl perl-dbd-mysql perl-time-hires graphviz

COPY --from=pprotein /go/src/app/pprotein /usr/local/bin/
COPY --from=pprotein /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=alp /go/bin/alp /usr/local/bin/
COPY --from=percona-toolkit /percona-toolkit-3.4.0/bin/* /usr/local/bin/

RUN mkdir -p /opt/pprotein
WORKDIR /opt/pprotein

ENTRYPOINT ["pprotein"]
