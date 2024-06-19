FROM golang:latest

# Install libvips and its dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libvips-dev \
    gcc \
    g++

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
COPY .env .

# Install the migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy entrypoint script
COPY scripts/run.sh .

# Build the Go app
RUN go build -ldflags "-s -w" -o main .
EXPOSE 8000

# Give execution permissions to the run script
RUN chmod +x run.sh

# Command to run the script
CMD ["./run.sh"]
