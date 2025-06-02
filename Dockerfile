FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ARG GIT_TAG
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-X main.version='"$GIT_TAG" -a -installsuffix cgo -o ./bin/server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
CMD [ "./server" ]
