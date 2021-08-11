FROM golang

WORKDIR /go/src/kinetur
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build ./cmd/web

ENTRYPOINT ["/go/bin/web"]
