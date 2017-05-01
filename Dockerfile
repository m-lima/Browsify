# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
RUN go version

# Install NPM
RUN curl -sL https://deb.nodesource.com/setup_6.x | bash - \
    && apt-get install -y nodejs

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/m-lima/browsify

# Build
RUN /go/src/github.com/m-lima/browsify/make.sh -o /opt/browsify

# Document the ports used by the image
EXPOSE 80
EXPOSE 443

# Run the server command by default when the container starts.
WORKDIR /opt/browsify
ENTRYPOINT /opt/browsify/browsify -c /opt/browsify/browsify.conf
