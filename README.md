# Simple Go Service 

This simple service was created to aid in setting up docker deployments

To start build an image and docker container to host the service create a Dockerfile:

```
FROM golang:1.8

EXPOSE 8080

WORKDIR /go/src/app/
COPY . .

RUN go get -d -v github.com/JamesWalter/simple-go-service
RUN go install -v github.com/JamesWalter/simple-go-service
CMD [ "simple-go-service" ]
```
Then build the image by executing this command in the same directory as the DockerFile
```
$ docker build --tag simple-go-service .
```
Then to run the go serivice will be listening on 8080, here 8080 is being mapped to 8888 on the docker host
```
$ docker run --detach --publish 8080:8888 --rm --name my-go_service simple-go-service 
```