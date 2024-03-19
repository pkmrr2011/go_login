# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21.7-alpine as builder

# Copy local code to the container image.
WORKDIR /go/src/app
COPY . .

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go build -v -o main .

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:latest  
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/main .
CMD ["./main"]  
