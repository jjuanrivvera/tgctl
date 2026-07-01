# tgctl — multi-stage build producing a tiny static image.
# Build:  docker build -t tgctl .
# Run:    docker run --rm -e TGCTL_TOKEN=123:ABC tgctl bot info

FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build

# Cache module downloads before copying the full source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static (CGO_ENABLED=0) binary. Version metadata matches the ldflags vars the Makefile
# and .goreleaser.yaml stamp into internal/version.
ARG VERSION=docker
ARG COMMIT=none
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
      -X github.com/jjuanrivvera/tgctl/internal/version.Version=${VERSION} \
      -X github.com/jjuanrivvera/tgctl/internal/version.Commit=${COMMIT} \
      -X github.com/jjuanrivvera/tgctl/internal/version.Date=${BUILD_DATE}" \
    -o /out/tgctl ./cmd/tgctl

# Distroless runtime: no shell, non-root, CA certs included for HTTPS to api.telegram.org.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/tgctl /usr/bin/tgctl
ENTRYPOINT ["/usr/bin/tgctl"]
CMD ["--help"]
