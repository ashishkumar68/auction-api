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

ARG TAG=latest
ARG FILE_UPLOADS_DIR=$HOME/uploads

ENV FILE_UPLOADS_DIR=$FILE_UPLOADS_DIR

RUN mkdir -p ${FILE_UPLOADS_DIR}
RUN chmod 755 -R ${FILE_UPLOADS_DIR}

EXPOSE 8081

RUN echo ${TAG}
#RUN go mod download
RUN go install github.com/ashishkumar68/auction-api@$TAG
RUN go install github.com/ashishkumar68/auction-api/auction-cli@$TAG

CMD ["auction-api"]