# build-stage
FROM golang:1.18.3-bullseye as build-stage
USER root

ENV TZ=Asia/Taipei

WORKDIR /
RUN mkdir build_space
WORKDIR /build_space
COPY . .
RUN go build -o toc-machine-trading ./cmd/app

# production-stage
FROM debian:bullseye as production-stage
USER root

ENV TZ=Asia/Taipei

WORKDIR /
RUN apt update -y && \
    apt install -y tzdata && \
    apt autoremove -y && \
    apt clean && \
    mkdir toc-machine-trading && \
    mkdir toc-machine-trading/data && \
    mkdir toc-machine-trading/migrations && \
    mkdir toc-machine-trading/configs && \
    mkdir toc-machine-trading/logs && \
    mkdir toc-machine-trading/scripts && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /toc-machine-trading

COPY --from=build-stage /build_space/toc-machine-trading ./toc-machine-trading
COPY --from=build-stage /build_space/data/holidays.json ./data/holidays.json
COPY --from=build-stage /build_space/migrations ./migrations/
COPY --from=build-stage /build_space/scripts/docker-entrypoint.sh ./scripts/docker-entrypoint.sh

ENTRYPOINT ["/toc-machine-trading/scripts/docker-entrypoint.sh"]
