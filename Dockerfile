FROM golang:1.26.4-alpine AS build

RUN apk add --no-cache ca-certificates
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cf-router ./cmd
RUN mkdir -p /out/data

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build --chown=65532:65532 /out/data /data
COPY --from=build /out/cf-router /cf-router

ENV ACCOUNT_PATH=/data/account.json WIREGUARD_PATH=/data/wg.json
USER 65532:65532
WORKDIR /data

ENTRYPOINT ["/cf-router"]