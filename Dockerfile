FROM golang:alpine AS builder

WORKDIR /build

COPY . /build/

RUN apk update && apk upgrade && apk add git && \
    go get -d

ARG TAG
ARG BUILDDATE
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=$BUILDDATE -X main.GitVersion=$TAG -extldflags \"-static\"" -o main .

FROM alpine:latest
LABEL maintainer="Andreas Peters <support@aventer.biz>"
LABEL org.opencontainers.image.title="mesos-compose"
LABEL org.opencontainers.image.description="ClusterD/Apache Mesos container orchestrator"
LABEL org.opencontainers.image.vendor="AVENTER UG (haftungsbeschränkt)"
LABEL org.opencontainers.image.source="https://github.com/AVENTER-UG/"

RUN apk add --no-cache ca-certificates
RUN apk update
RUN apk upgrade
RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/main /app/

EXPOSE 10000

WORKDIR "/app"

CMD ["./main"]
