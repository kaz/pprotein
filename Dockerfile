# --------------------------------------------------

FROM alpine AS pprotein

RUN apk add go npm make

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN make view/dist
RUN make pprotein

# --------------------------------------------------

FROM alpine AS kataribe

RUN apk add go

ENV GOPATH /go
RUN go get github.com/matsuu/kataribe

# --------------------------------------------------

FROM alpine AS percona-toolkit

RUN wget https://downloads.percona.com/downloads/percona-toolkit/3.3.1/binary/tarball/percona-toolkit-3.3.1_x86_64.tar.gz
RUN tar zxvf percona-toolkit-3.3.1_x86_64.tar.gz

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache bash perl perl-dbd-mysql perl-time-hires graphviz

COPY --from=pprotein /go/src/app/pprotein /usr/local/bin/
COPY --from=kataribe /go/bin/kataribe /usr/local/bin/
COPY --from=percona-toolkit /percona-toolkit-3.3.1/bin/* /usr/local/bin/

RUN mkdir -p /opt/pprotein
WORKDIR /opt/pprotein

ENTRYPOINT ["pprotein"]
