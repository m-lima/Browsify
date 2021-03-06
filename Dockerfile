## Backend
# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
RUN go version

WORKDIR /go/src/github.com/m-lima/browsify

# Copy the local package files to the container's workspace.
COPY *.go /go/src/github.com/m-lima/browsify/
COPY auther /go/src/github.com/m-lima/browsify/auther

# Build
RUN go get && go install

## Main
FROM golang

COPY --from=browsify-web /web/build /opt/browsify/web
COPY --from=0 /go/bin/browsify /opt/browsify/.
COPY secrets/* /opt/browsify/
COPY *.conf /opt/browsify/

# Document the ports used by the image
EXPOSE 80

# Run the server command by default when the container starts.
WORKDIR /opt/browsify
CMD ["./browsify"]
