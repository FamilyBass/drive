# Build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG GOOS=linux
ARG GOARCH=arm64
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o /out/server ./

# Final stage
FROM scratch
COPY --from=builder /out/server /server
EXPOSE 8080
VOLUME ["/data","/app/media"]
ENTRYPOINT ["/server"]

