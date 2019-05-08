FROM golang:stretch

WORKDIR $GOPATH/src/github.com/yalotso/thumbnail

ENV GO111MODULE=on
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go test -cover ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

EXPOSE 8080
ENTRYPOINT ["./app"]