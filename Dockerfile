FROM golang:1.25.5-alpine AS build
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Add dependencies to compile sqlite
RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace
# Setup orchestrion
RUN go install github.com/DataDog/orchestrion@latest

# Copy dependencies files
COPY go.mod go.sum ./
# Downloads dependencies
RUN go mod download

# Copy source code
COPY ./*.go /workspace/
# Pins orchestrion dependencies
RUN orchestrion pin
# Builds a static trimmed binary with orchestrion injectione
RUN go build -ldflags='-s -w -extldflags "-static"' -toolexec="orchestrion toolexec" -o "moe-count"

FROM scratch AS deploy

WORKDIR /app/
COPY --from=build /workspace/moe-count /usr/local/bin/moe-count
COPY ./static /app/static
COPY ./template /app/template

ENTRYPOINT [ "moe-count" ]
