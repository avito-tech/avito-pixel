FROM golang:1.20

WORKDIR /avito-pixel
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o /avito-pixel-server ./example/service/
EXPOSE 3000
CMD ["/avito-pixel-server"]
