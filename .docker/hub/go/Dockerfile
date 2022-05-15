FROM golang:1.18-alpine

RUN runtimeDeps=" \
        curl \
        vim \
        nano \
        htop \
        net-tools \
    " \
    && apk update && apk add -y $runtimeDeps

ARG TAG=latest
ARG UPLOADS_DIR=${HOME}/uploads

ENV UPLOADS_DIR=$UPLOADS_DIR

RUN mkdir -p ${UPLOADS_DIR}
RUN chmod 755 -R ${UPLOADS_DIR}

#ENV GOPATH='/go-workspace'
#ENV GOBIN=${GOPATH}/bin
#ENV GOSRC=${GOPATH}/src
#ENV PROJECTDIR=${GOSRC}/auction-api
## add go bin to path, this is required to run go binaries.
#ENV PATH=$PATH:${GOBIN}

EXPOSE 8081

#RUN mkdir -p ${GOPATH}/bin/
#RUN mkdir -p ${GOPATH}/pkg/
#RUN mkdir -p ${GOPATH}/src/

#COPY . ${PROJECTDIR}
#WORKDIR ${PROJECTDIR}

#RUN go mod download
RUN go install github.com/ashishkumar68/auction-api@${TAG}
RUN go install github.com/ashishkumar68/auction-api/auction-cli@${TAG}

CMD ["auction-api"]