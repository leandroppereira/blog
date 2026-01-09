# Build
FROM golang:1.22 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/blog-sim ./main.go

# Run (distroless)
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/blog-sim /blog
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/blog]
