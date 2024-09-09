FROM golang:1.23.1-alpine AS builder

ARG APP_VNUM=demo
ARG APP_NAME=component

# Move to working directory (/build).
WORKDIR /build

RUN apk add --no-cache upx 
# RUN apk add --no-cache make gcc g++

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environmet variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.version=$APP_VNUM" -o $APP_NAME main.go
RUN upx --best --lzma $APP_NAME

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/component", "/"]

# Export necessary port.
EXPOSE 5000

# Command to run when starting the container.
ENTRYPOINT ["/component"]
