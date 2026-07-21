# Stage 1: build the static binary.
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /out/plug-scaffold \
    ./cmd/plug-scaffold

# Stage 2: minimal runtime image. alpine (not scratch) gives a shell for
# debugging and matches the sibling neospec image shape.
FROM alpine:3.20

RUN apk add --no-cache git

COPY --from=builder /out/plug-scaffold /usr/local/bin/plug-scaffold

ENTRYPOINT ["/usr/local/bin/plug-scaffold"]
CMD ["--help"]
