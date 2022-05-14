FROM golang:1.18

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

ENV GOPATH='/go-workspace'
ENV GOBIN=${GOPATH}/bin
ENV GOSRC=${GOPATH}/src
ENV PROJECTDIR=${GOSRC}/auction-api
# add go bin to path, this is required to run go binaries.
ENV PATH=$PATH:${GOBIN}

EXPOSE 8081

RUN mkdir -p ${GOPATH}/bin/
RUN mkdir -p ${GOPATH}/pkg/
RUN mkdir -p ${GOPATH}/src/

RUN ls -al

COPY . ${PROJECTDIR}
WORKDIR ${PROJECTDIR}

RUN ls -al

COPY ./init/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

RUN go mod download
RUN go install
RUN go install github.com/ashishkumar68/auction-api/auction-cli@latest

CMD ["./auction-api"]