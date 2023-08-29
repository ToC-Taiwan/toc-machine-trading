# build-stage
FROM golang:1.20.7-bullseye as build-stage
USER root

ENV TZ=Asia/Taipei

WORKDIR /
RUN mkdir build_space
WORKDIR /build_space
COPY . .
RUN CGO_ENABLED=0 go build -o toc-machine-trading ./cmd/app

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

ENTRYPOINT ["/toc-machine-trading/toc-machine-trading"]
