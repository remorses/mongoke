# Start from the latest golang base image
FROM golang:latest as builder

# RUN apk --no-cache add ca-certificates build-base

# Set the Current Working Directory inside the container
WORKDIR /cwd

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY src  ./src
COPY cmd  ./cmd



# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go
# RUN ldd main

######## Start a new stage from scratch #######
FROM alpine:latest


# RUN apk --no-cache add ca-certificates # libc6-compat

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /cwd/main .
COPY ./web-ui/out ./web-ui/

RUN ls -l
RUN ls -l ./web-ui

# RUN ldd ./main

# Expose port 8080 to the outside world
EXPOSE 8080

ENV WEB_UI_ASSETS=./web-ui

# Command to run the executable
ENTRYPOINT ["./main"]