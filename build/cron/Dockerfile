FROM golang:1.19-alpine3.16 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/cron -v ./cmd/cron

FROM alpine:3.16

# Copy the binary from the builder image
COPY --from=builder /bin/cron .

CMD [ "sh", "-c", "./main" ]
