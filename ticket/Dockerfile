FROM golang:1.19.13-alpine3.18 as builder
WORKDIR /app
RUN apk add git gcc libc-dev make
ENV GO111MODULE on
ENV CGO_ENABLED=1
ENV GOOS=linux
COPY go.mod .
COPY go.sum .
COPY . .
RUN apk update
RUN make build

FROM alpine:3.18 as release
RUN apk add --no-cache --update ca-certificates
COPY --from=builder /app/main /app/cmd/
COPY --from=builder /app/asset /app/asset
RUN chmod +x /app/cmd/main
WORKDIR /app
CMD ["cmd/main", "serve-rest"]
