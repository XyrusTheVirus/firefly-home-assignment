#builder: compile stage
FROM golang:latest as builder
WORKDIR /firefly/app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy sources
COPY . ./
RUN mv .env.example .env

# Set GOBIN to a safe place for the binary
ENV GOBIN=/firefly/app/bin
RUN mkdir -p $GOBIN

# Install the binary
RUN go install ./cmd

# ---------------------------------------
# Development stage (with air)
FROM golang:latest as dev-builder
WORKDIR /firefly/app
COPY --from=builder /firefly/app/.env .env
COPY . ./
RUN go install github.com/air-verse/air@latest
ENV AIR_CONFIG=/firefly/app/.air.toml
ENTRYPOINT ["air"]


# ---------------------------------------
# Production stage
FROM debian:13-slim as prod-builder
WORKDIR /firefly/app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
# Copy installed binary from $GOBIN
COPY --from=builder /firefly/app/bin/cmd ./firefly

# Copy runtime files
COPY --from=builder /firefly/app/.env .env
COPY --from=builder /firefly/app/endg-urls endg-urls
COPY --from=builder /firefly/app/words.txt words.txt

CMD ["./firefly"]