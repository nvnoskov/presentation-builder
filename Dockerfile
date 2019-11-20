# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Nikolay Noskov <me@noskov.dev>"

WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
ARG GO111MODULE=on
ARG CGO_ENABLED=1

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go get github.com/shiyanhui/hero/hero && go get golang.org/x/tools/cmd/goimports


# RUN go mod vendor
COPY . .
# RUN apt-get update && apt-get install -y libmupdf-dev libsqlite3-dev

RUN hero -source=./templates && go build -ldflags "-s -w" -o app

# Copy the source from the current directory to the Working Directory inside the container
EXPOSE 8000

# Build the Go app
# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .
# RUN go build -a -installsuffix cgo -o main .

FROM ubuntu:eoan
COPY --from=builder /app /
CMD ["/app"]


# FROM alpine:latest  

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/main .

# Expose port 8000 to the outside world


# Command to run the executable
# CMD ["./main"] 