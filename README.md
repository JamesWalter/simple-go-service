# Simple Go Service 

This repository was created to help learn about docker, go, nginx, and kubernetes.

# Go Web App 

The go app itself is a simple web app that will display the hostname (container id), the time the web app was started, and the number of visits

For example
```
Hello from container "45d34b1449bd"
I have been visted 38 times
I was born "Thu Mar  8 22:39:24 UTC 2018" 
```

# Build Docker Image

To get the app up and running in a container we need to build a Docker image. Start by creating a new directory. e.g.

```
mkdir $HOME\my-docker-build
touch $HOME\my-docker-build\Dockerfile
```
Open the new Dockerfile in your favorite text editor (e.g. vim) and then write instructions to build an image.

### Dockerfile
```
FROM golang:1.8

EXPOSE 8080

WORKDIR /go/src/app/
COPY . .

RUN go get -d -v github.com/JamesWalter/simple-go-service
RUN go install -v github.com/JamesWalter/simple-go-service
CMD [ "simple-go-service" ]
```

### Build Command
To build a Docker container image the Dockerfile will be executed

The command below is the simlest the build can get, this command must be executed from the directory of the dockerfile
```
$ docker build --tag simple-go-service .
```
* `--tag simple-go-service` will tag the resulting image. When you build the image the last message will be `Successfully tagged simple-go-service:latest`. The message indicates the image name as `simple-go-service` and the tag as `latest`, the tag represents a specfic version. Docker by default applies the `latest` tag. Generally when tagging with intent to deploy you'll want to give a version to tag the image along with the name for instance `--tag my-app:1.0`.

# Running the Container
To run the go web app the docker run command will be used. 
```
$ docker run --detach --publish 8888:8080 --rm --name go-1 simple-go-service 
```
Here is an explanation of what is going on with the command

* `--detach`
* `--publish 8888:8080`
* `--rm`
* `--name go-1`

### Check running container

### Access web app
* browser
* curl

# Running multiple containers
Now run mulitple containers each running its own instance of the go web app

# Nginx

### Nginx Load Balancer config
Create `nginx.conf` file
```
worker_processes 4;

events { worker_connections 1024; }

http {
    sendfile on;

    upstream app_servers {
        server go-1:8080;
        server go-2:8080;
        server go-3:8080;
    }
    
    server {
	    listen 8888;
	    location / {
            proxy_pass         http://app_servers;
            proxy_redirect     off;
            proxy_set_header   Host $host;
            proxy_set_header   X-Real-IP $remote_addr;
            proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header   X-Forwarded-Host $server_name;
	    }
    }
}

```

### Nginx Dockerfile
```
FROM nginx
COPY nginx.conf /etc/nginx/nginx.conf
```

### Build Image


### Set up network 
```
$ docker network create test-net
```
```
$ docker run --detach --network test-net --name go-1 simple-go-service
$ docker run --detach  --network test-net --name go-2 simple-go-service
$ docker run --detach  --network test-net --name go-3 simple-go-service
```

### Run Nginx
```
$ docker run --detach --rm --network test-net -p 8888:8888 --name proxy go-service-proxy
````


