FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/google-scholar-mcp ./cmd/google-scholar-mcp

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/google-scholar-mcp /usr/local/bin/google-scholar-mcp

ENTRYPOINT ["/usr/local/bin/google-scholar-mcp", "stdio"]
