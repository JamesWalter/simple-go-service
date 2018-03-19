FROM golang:alpine

EXPOSE 8080

WORKDIR /go/src/app
COPY ./main.go .

RUN go get -v app
RUN go install -v app
ENTRYPOINT [ "/go/bin/app" ]