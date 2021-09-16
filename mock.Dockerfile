# --------------------------------------------------

FROM alpine AS pprotein

RUN apk add go npm make

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN make build

# --------------------------------------------------

FROM alpine AS mock

RUN apk add go

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/app
COPY . .

RUN go build ./cli/pprotein-mock

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache mysql mysql-client nginx supervisor
RUN mysql_install_db --datadir=/var/lib/mysql --basedir=/usr --user=root
RUN mkdir /var/log/mysql

COPY --from=pprotein /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=mock /go/src/app/pprotein-mock /usr/local/bin/

COPY mock/supervisord.ini /etc/supervisor.d/
COPY mock/mysqld.cnf /etc/my.cnf.d/
COPY mock/nginx.conf /etc/nginx/

ENV DSN "root@unix(/var/run/mysqld/mysqld.sock)/"
ENV REQUEST_HOST "127.0.0.1:80"

ENTRYPOINT ["supervisord"]
