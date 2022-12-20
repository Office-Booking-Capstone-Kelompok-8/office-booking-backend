FROM golang:1.19-alpine3.16 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o /bin/main -v ./cmd/main

FROM alpine:3.16

# Copy the binary from the builder image
COPY --from=builder /bin/main .

#  Install tzdata for timezone support
RUN apk add --no-cache tzdata

# Install and use dumb-init to give fiber a PID to inherit
RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD [ "sh", "-c", "./main" ]
