ARG GOLANG_VERSION=1.25-alpine3.21
ARG ALPINE_VERSION=3.21

FROM golang:${GOLANG_VERSION} AS deps

WORKDIR /app

COPY ./ ./

ENV GO111MODULE=on

WORKDIR /app

ENV CGO_ENABLED=0
ARG ARTIFACT_VERSION

RUN apk add --no-cache make

RUN make build 

FROM deps AS build

FROM alpine:${ALPINE_VERSION} AS runtime

WORKDIR /app

COPY --from=build /app/bin /app

RUN apk update \
    && apk add --no-cache --upgrade \
        bash \
        ca-certificates \
        curl \
        tzdata \
    && update-ca-certificates \
    && echo 'Etc/UTC' > /etc/timezone \
    && adduser --disabled-password --home /app --gecos '' gouser \
    && chown -R gouser /app

ENV TZ     :/etc/localtime
ENV LANG   en_US.utf8
ENV LC_ALL en_US.UTF-8

USER gouser

ENTRYPOINT ["/app/auth"]
