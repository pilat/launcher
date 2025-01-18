FROM golang:1.23.4 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /out/launcher .

FROM scratch

LABEL org.opencontainers.image.title="Launcher" \
      org.opencontainers.image.description="A lightweight launcher application built with Go" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.authors="Vladimir Urushev <vkurushev@gmail.com>" \
      org.opencontainers.image.source="https://github.com/pilat/launcher"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/launcher /usr/local/bin/launcher

CMD ["launcher"]
