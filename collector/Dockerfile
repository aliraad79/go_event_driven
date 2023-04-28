FROM golang:1.20

WORKDIR /collector

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./go_collector

EXPOSE 8080

RUN chmod +x ./go_collector

# Run
ENV GIN_MODE=release
CMD ["./go_collector"]