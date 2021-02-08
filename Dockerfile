FROM golang:1.15-alpine AS builder
ADD . /dictiobot
WORKDIR /dictiobot
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o dictiobot cmd/dictiobot/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
COPY --from=builder /dictiobot ./
RUN chmod +x dictiobot
CMD ./dictiobot