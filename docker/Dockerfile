FROM ubuntu:18.04 as builder

# --progress string

RUN sed -i "s/archive.ubuntu.com/mirror.kakao.com/g" /etc/apt/sources.list
RUN apt-get update 
RUN apt-get install -y wget libzip-dev pkg-config

RUN wget https://go.dev/dl/go1.17.3.linux-amd64.tar.gz && \
    tar -zxf go1.17.3.linux-amd64.tar.gz && \
    mv go /usr/lib && \ 
    update-alternatives --install /usr/bin/go go /usr/lib/go/bin/go 40 && \
    update-alternatives --query go

RUN mkdir -p /build/
COPY . /build

ENV PKG_CONFIG_PATH /build/zlib/linux

RUN go version
RUN cd /build/extractor && go build && \
    cd /build/routed && go build


# 
FROM ubuntu:18.04 as runstage
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Seoul

RUN sed -i "s/archive.ubuntu.com/mirror.kakao.com/g" /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y wget

COPY --from=builder /usr/lib/go /usr/lib
COPY --from=builder /build/extractor/extractor /usr/local/bin
COPY --from=builder /build/routed/routed /usr/local/bin

RUN extractor; exit 0 
RUN routed; exit 0
