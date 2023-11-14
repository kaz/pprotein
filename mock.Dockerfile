# --------------------------------------------------

FROM golang:alpine AS agent

WORKDIR $GOPATH/src/app
COPY . .

RUN go build ./cli/pprotein-agent

# --------------------------------------------------

FROM golang:alpine AS mock

WORKDIR $GOPATH/src/app
COPY . .

RUN go build ./cli/pprotein-mock

# --------------------------------------------------

FROM alpine/git AS repo

WORKDIR /opt
RUN git clone https://github.com/kaz/pprotein.git

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache mysql mysql-client nginx supervisor
RUN mysql_install_db --datadir=/var/lib/mysql --basedir=/usr --user=root
RUN mkdir /var/log/mysql

COPY --from=agent /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=mock /go/src/app/pprotein-mock /usr/local/bin/
COPY --from=repo /opt/pprotein /opt/pprotein

COPY mock/supervisord.ini /etc/supervisor.d/
COPY mock/mysqld.cnf /etc/my.cnf.d/
COPY mock/nginx.conf /etc/nginx/

ENV DSN "root@unix(/var/run/mysqld/mysqld.sock)/"
ENV REQUEST_HOST "127.0.0.1:80"
ENV PPROTEIN_GIT_REPOSITORY "/opt/pprotein"

ENTRYPOINT ["supervisord"]
