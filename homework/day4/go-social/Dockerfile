# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22.3 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY . .

  RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gosocial

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  WORKDIR /

  COPY --from=build-stage /bin/gosocial /bin/gosocial

  EXPOSE 8080

  ENTRYPOINT ["/bin/gosocial"]