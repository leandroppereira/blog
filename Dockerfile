# Build
FROM golang:1.22 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /blog ./main.go

# Runtime com shell (compatível com OpenShift)
FROM registry.access.redhat.com/ubi9/ubi-micro
COPY --from=builder /blog /blog
EXPOSE 8080
USER 1001
ENTRYPOINT ["/blog"]