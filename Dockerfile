FROM golang:1.21-alpine

WORKDIR /app

# Install libpcap-dev and libpcap for gopacket dependencies
# Also, install build-base (or gcc and musl-dev) if you encounter issues with cgo
RUN apk update && apk add --no-cache libpcap-dev libpcap gcc musl-dev

# Install air for live reloading
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]