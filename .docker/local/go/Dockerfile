FROM golang:1.18

ARG AIR_FILE_NAME=".air.toml"

RUN buildDeps=" \
        lsb-release \
    " \
    runtimeDeps=" \
        curl \
        git \
        vim \
        nano \
        htop \
        net-tools \
    " \
    && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y $buildDeps $runtimeDeps

RUN wget https://dev.mysql.com/get/mysql-apt-config_0.8.22-1_all.deb
RUN dpkg -i mysql-apt-config_0.8.22-1_all.deb
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt install -y mariadb-client

ENV GOPATH='/go-workspace'
ENV GOBIN=${GOPATH}/bin
ENV GOSRC=${GOPATH}/src
ENV PROJECTDIR=${GOSRC}/auction-api
# add go bin to path, this is required to run go binaries.
ENV PATH=$PATH:${GOBIN}
ENV AIR_FILE_NAME ${AIR_FILE_NAME}

RUN mkdir -p ${GOPATH}/bin/
RUN mkdir -p ${GOPATH}/pkg/
RUN mkdir -p ${GOPATH}/src/
# install wire and other commands.
RUN go install github.com/google/wire/cmd/wire@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/spf13/cobra-cli@latest
RUN go install github.com/ashishkumar68/auction-api/auction-cli@latest

EXPOSE 8081

COPY init/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

WORKDIR ${PROJECTDIR}
CMD air -c ${AIR_FILE_NAME}