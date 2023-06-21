FROM docker.io/library/golang:latest  as builder

WORKDIR /

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o server

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /server /server
COPY --from=builder /static /static
COPY --from=builder /templates /templates
ENV GIN_MODE release

CMD ["/server"]
