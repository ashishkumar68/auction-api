FROM golang:1.17

RUN buildDeps=" \
        lsb-release \
    " \
    runtimeDeps=" \
        curl \
        git \
        vim \
        nano \
        htop \
    " \
    && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y $buildDeps $runtimeDeps

RUN wget https://dev.mysql.com/get/mysql-apt-config_0.8.22-1_all.deb
RUN dpkg -i mysql-apt-config_0.8.22-1_all.deb
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt install -y mariadb-client

ENV GOPATH='/go-workspace'
ENV GOBIN=${GOPATH}/bin
ENV GOSRC=${GOPATH}/src
ENV PROJECTDIR=${GOSRC}/github.com/ashishkumar68/auction-api
# add go bin to path, this is required to run go binaries.
ENV PATH=$PATH:${GOBIN}

RUN mkdir -p ${GOPATH}/bin/
RUN mkdir -p ${GOPATH}/pkg/
RUN mkdir -p ${GOPATH}/src/
# install wire
RUN go install github.com/google/wire/cmd/wire@latest

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

EXPOSE 8081

COPY init/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

WORKDIR ${PROJECTDIR}
CMD ["air"]