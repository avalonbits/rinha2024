FROM golang:1.22 as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -tags 'fts5,osusergo,netgo,static' --ldflags '-linkmode external -extldflags "-static"' -o /app/rinha ./cmd/rinha

EXPOSE 1323

# Run on container startup.
CMD ["/app/rinha"]
