FROM golang:1.9.1

ENV GOPATH=/go:/go/src/github.com/morrah77/rates
WORKDIR /go/src/github.com/morrah77/rates
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN cd src/rates/ && dep ensure && cd ../..
RUN go install rates/main
CMD ./bin/main --listen-addr :8080
