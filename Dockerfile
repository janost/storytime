FROM golang:1.21-alpine as gobuild
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch
WORKDIR /app
COPY --from=gobuild /app/storytime .

ENTRYPOINT ["/app/storytime"]
