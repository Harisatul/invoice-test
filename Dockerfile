# Use an official Golang runtime as the base image  
FROM golang:1.23-alpine3.20 AS BuildStage
  
# Set the working directory inside the container  
WORKDIR /go/src/invoice-test

COPY . .

# Build the Go app
RUN go build ./main.go


# Deploy Stage
FROM alpine:latest

COPY app.yaml /

# Copy the built executable from BuildStage to the root directory
COPY --from=BuildStage /go/src/invoice-test /

ENTRYPOINT ["./main", "serve-http"]
