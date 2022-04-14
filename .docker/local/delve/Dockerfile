FROM golang:1.18 as build-env

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

RUN mkdir -p ${GOPATH}/bin/
RUN mkdir -p ${GOPATH}/pkg/
RUN mkdir -p ${GOPATH}/src/

ENV GOBIN=${GOPATH}/bin
ENV GOSRC=${GOPATH}/src
ENV PROJECTDIR=${GOSRC}/auction-api
# add go bin to path, this is required to run go binaries.
ENV PATH=$PATH:${GOBIN}
ENV CGO_ENABLED 0
ENV DELVE_PORT 40000

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8081 40000

COPY init/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

WORKDIR ${PROJECTDIR}

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--log=true", "--log-output=debugger,debuglineerr,gdbwire,lldbout,rpc", "exec", "/go-workspace/bin/server"]
#CMD ["go", "run", "server.go"]