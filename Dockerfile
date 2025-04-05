# Step 1: Use Go version 1.19 for building the app (valid version)
FROM golang:1.23.0 AS build

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all the files and build the application
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final aşama
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/config_files ./config_files
COPY --from=build /app/specific_configs ./specific_configs

COPY --from=build /app/main .

EXPOSE 8000
CMD ["/root/main"]