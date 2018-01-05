FROM golang:latest

# install depedencies
RUN go get github.com/pressly/goose/cmd/goose

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/cryplio/rest-api

# Install api binary globally within container
RUN cd /go/src/github.com/cryplio/rest-api && make install

# Set binary as entrypoint
CMD /go/bin/cryplio-api

EXPOSE 5002
