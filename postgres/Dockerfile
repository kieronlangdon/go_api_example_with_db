FROM golang:alpine3.12 AS builder
RUN apk add git
ENV CGO_ENABLED=0

########
# Prep
########

# add the source
COPY . /go/src/postgres
WORKDIR /go/src/postgres

########
# Build Go Wrapper
########

# Get the dependencies
RUN go get -v -t -d ./...

# Build the package tests
RUN if [ ! -d "package.test" ] ; then echo Test package does not exist ; else rm-rf package.test ; fi
RUN go test -c . -o package.test
RUN go test -v -run Unit

# Build the go app
RUN go build -o ./postgres *.go

########
# Package into runtime image (smaller image size)
########
FROM alpine as release

# Copy the executable from the builder image
COPY --from=builder /go/src/postgres .

CMD ["/postgres"]