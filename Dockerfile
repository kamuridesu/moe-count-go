FROM golang:1.21-alpine AS build
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . /workspace/
RUN go build -ldflags='-s -w -extldflags "-static"' -o "moe-count"

FROM scratch AS deploy
ENV GIN_MODE=release

WORKDIR /app/
COPY --from=build /workspace/moe-count /usr/local/bin/moe-count
COPY ./static /app/static
COPY ./template /app/template

ENTRYPOINT [ "moe-count" ]
