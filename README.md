# Simple Go Service 

This simple service was created to aid in setting up docker deployments

To start a docker container to host the service create a Dockerfile:

```
FROM golang:1.8

EXPOSE 8080

WORKDIR /go/src/app/
COPY . .

RUN go get -d -v /github.com/JamesWalter/simple-go-service
RUN go install -v /github.com/JamesWalter/simple-go-service
CMD [ "simple-go-service" ]
```

To build the docker image
```
$ docker build -t test-go.
```
Then to run 
Using auto port mapping
```
$ docker run -it --rm --name my-test-go test-go -P 
```
Specify your own port mapping
```
$ docker run -it --rm --name my-test-go test-go -p 8080:8080
```