FROM golang:1.20

# Set working dir
WORKDIR /app

# Copy all go files
COPY src .

# Run the tests
RUN go test -v ./...

# Build the package
RUN go build -o wbtest ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/app/wbtest"]