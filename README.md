# Run the binary locally
./wbtest

# Build Docker image (run the tests/gosec, build binary, expose port)
docker build -f Dockerfile -t wbtest:latest .

# Run Docker image in backround and forward port 
docker run -d -p 8080:8080 wbtest

# Save the user
curl -v localhost:8080/save -d '{"id":"123", "name":"Some Name", "email":"email@email.com","date_of_birth":"2020-01-01T12:12:34+00:00"}'

# Get the user data
curl -v localhost:8080/123