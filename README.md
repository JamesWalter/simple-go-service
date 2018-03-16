# Simple Go Service 

This repository was created to help learn about [Docker](https://www.docker.com/), [Go](https://en.wikipedia.org/wiki/Go_%28programming_language%29), [Nginx](https://en.wikipedia.org/wiki/Nginx), and [Kubernetes](https://en.wikipedia.org/wiki/Kubernetes).

The tutorial below can be run in a free in bowser Alpine Linux Virtual Machine. [Play With Docker](https://labs.play-with-docker.com/)

For more information about docker:
+ [Docker offical website](https://www.docker.com/) 
+ [Docker tutorials using play with docker](https://training.play-with-docker.com/)
+ [Docker Get Started](https://docs.docker.com/get-started/)

The tutorial below will cover
+ Building a Golang Docker image that contains simple-go-service
+ Running containers 
+ Buidling a nginx reverse proxy that will manage connections to mutiple simple-go-service containers
+ (Todo) Orchastrate nginx and go containers with Kubernetes

# Go Web App 

The Go app itself is a simple web app that will return the hostname (container id), the time the web app was started, and the number of visits.

For example
```
Hello from container "45d34b1449bd"
I have been visted 38 times
I was born "Thu Mar  8 22:39:24 UTC 2018" 
```

# Build Docker Image

To get the app up and running in a container we need to build a [Docker image](https://docs.docker.com/glossary/?term=image). 

Start by creating a new directory called `my-first-image`. Then create a new file named `Dockerfile` in the directory.

If you are using  play with docker you will be using the command line. The commands denoted with $ are command line commands
```
$ mkdir "$HOME"/my-first-image
$ touch "$HOME"/my-first-image/Dockerfile
```

### Dockerfile

A `Dockerfile` is the cookbook for a docker image. It contains all of the commands, in order, necessary to construct a docker image. 

Open the new `Dockerfile` in your favorite text editor (e.g. [vim](https://www.howtoforge.com/vim-basics)).
```
$ cd "$HOME"/my-first-image
$ vim Dockerfile
```

Here are the contents of the enire `Dockerfile`
```
FROM golang:1.8

EXPOSE 8080

RUN go get -d -v github.com/JamesWalter/simple-go-service
RUN go install -v github.com/JamesWalter/simple-go-service
CMD [ "simple-go-service" ]
```

Here is a brief explanation of each line in the `Dockerfile`
+ `FROM golang:1.8` This image is to be an extension of the Golang image version 1.8 on [Docker Hub](https://hub.docker.com/_/golang/)
+ `EXPOSE 8080` Lets Docker know that the container will listen on port 8080. The Go web app is configured to listen on port 8080.
+ `WORKDIR /go/src/app/` Set the present working directory within the container to /go/src/app'.
+ `COPY . .` copy the contents of the present working directory to the present working directory in the the container
+ `RUN go get -d -v github.com/JamesWalter/simple-go-service` Run command to get the contents of the git hub repository within the container (this is a Golang specific comman)
+ `RUN go install -v github.com/JamesWalte/simple-go-service` Install the go application within the container
+ `CMD [ "simple-go-service" ]` Execute command `simple-go-service` within the container

See [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/) for additional information

### Build Command
To build a Docker container image the `Dockerfile` will be executed by issuing a `$ docker build` command.

The command below is the simple. This command must be executed from the directory of the `Dockerfile`
```
$ docker build --tag simple-go-service .
```
* `--tag simple-go-service` will tag the resulting image. By default Docker tags the images with version `latest`

Find out more about the build command by running `$ docker build --help` or by visiting [Docker Build Reference](https://docs.docker.com/engine/reference/commandline/build/)

# Running the Container
To run the go web app the `docker run` command will be used. 
```
$ docker run --detach --publish 8888:8080 --rm --name go-1 simple-go-service 
```
Here is an explanation of what is going on with the command

* `--detach` Run the container in the background
* `--publish 8888:8080` Map the listening port of the container (8080) to the hostmachine's port (8888).
* `--rm` Delete the container once execution has finished
* `--name go-1` Name the container go-1

Additional options can be found by running `$ docker run --help` or by visiting [Docker Run Reference](https://docs.docker.com/engine/reference/run/)

### Check running container
To check that the container is running execute `docker ps` you should see a display similiar to the one below. This shows the active containers
```
$ docker ps
CONTAINER ID        IMAGE               COMMAND               CREATED             STATUS              PORTS                    NAMES
14d5c6c19f13        simple-go-service   "simple-go-service"   3 minutes ago       Up 3 minutes        0.0.0.0:8888->8080/tcp   go-1
```

### Access web app
To see the web app in action visit it in one of two ways

* From a web browser visit `http://localhost:8888`
* From the command line `$ curl http://127.0.0.1:8888`

Here is what a curl command should return.
```
$ curl http://127.0.0.1:8888
Hello from container "14d5c6c19f13"
I have been visted 1 times
I was born "Mon Mar 12 21:42:35 UTC 2018"
```

# Running multiple containers
Now run mulitple containers each running its own instance of the Go web app. Do this by executing the `docker run` command using different names and host port mapping for each container

```
$ docker run --detach --publish 8889:8080 --rm --name go-2 simple-go-service 
$ docker run --detach --publish 9000:8080 --rm --name go-3 simple-go-service 
```

Now execute `docker ps` and all the new containers should now show up as running
```
$ docker ps
CONTAINER ID        IMAGE               COMMAND               CREATED             STATUS              PORTS                    NAMES
b2566f5e00d6        simple-go-service   "simple-go-service"   4 seconds ago       Up 3 seconds        0.0.0.0:9000->8080/tcp   go-3
2adf608a9e04        simple-go-service   "simple-go-service"   12 seconds ago      Up 12 seconds       0.0.0.0:8889->8080/tcp   go-2
14d5c6c19f13        simple-go-service   "simple-go-service"   15 minutes ago      Up 15 minutes       0.0.0.0:8888->8080/tcp   go-1
```

Try visiting each different host port, paying attention to the container id's.

# Nginx

Stepping complexity up a bit, its time build a [Nginx reverse proxy server](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/). The server will be the single point of access for the Go web app. Additionally the server will distribute accesses amongst multiple Go containers.

Before building the reverse proxy server do some house keeping. Stop all running instances of the go web app. Use `docker ps` and `docker stop <container_id>` or `docker stop <container_name>`. For example 
```
$ docker stop go-1
```

Now create a new directory for the proxy server build, add a Dockerfile, and a file named `nginx.conf`.
```
$ mkdir "$HOME"/nginx
$ touch "$HOME"/nginx/Dockerfile
$ touch "$HOME"/nginx/nginx.conf

```

### Nginx Reverse Proxy config
The `nginx.conf` file will contain the configuration of the nginx server that will be deployed in the Docker container.

Navigate to and open the nginx.conf file.

```
$ cd "$HOME"/nginx
$ vim nginx.conf
```

Add the following configuration to the file
```
worker_processes 1;

events{ worker_connections 1028; }

http {

    upstream app_servers {
        server go-1:8080;
        server go-2:8080;
        server go-3:8080;
    }
    
    server {
	    listen 8888;
	    location / {
            proxy_pass         http://app_servers;
            proxy_set_header   Host $host;
            proxy_set_header   X-Real-IP $remote_addr;
            proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header   X-Forwarded-Host $server_name;
	    }
    }
}

```

The deatails surrounding the configuration won't be discussed here but a good starting point for learning more is the [Nginx Beginners Guide](http://nginx.org/en/docs/beginners_guide.html)

### Nginx Dockerfile
The Dockerfile for the Nginx reverse proxy is simple, it extends the nginx image from [Docker Hub](https://hub.docker.com/_/nginx/) and copies the configuration file that was created into the image.

Open the Dockerfile that was created in the nginx directory.
```
$ cd "$HOME"/nginx
$ vim Dockerfile
```

Enter the commands below and save
```
FROM nginx
COPY nginx.conf /etc/nginx/nginx.conf
```

### Build Image
Building the Nginx reverse proxy image is simple. From within the nginx directory execute the `docker build` command
```
$ cd "$HOME"/nginx
$ docker build --tag nginx-reverse-proxy .
```

### Set up network 
Before starting up simple-go-service containers and the nginx-proxy container, set up a Docker network. This network will allow connections between containers without having to expose ports on the host machine.
```
$ docker network create test-net
```

### Run multiple Go contianers
With the network `test-net` setup. Run 3 simple-go-service containers, include the `--network test-net` option. Notice that `--publish' has been removed. The listening ports for each container are only available to other containers within Docker.

```
$ docker run --detach --network test-net --name go-1 simple-go-service
$ docker run --detach --network test-net --name go-2 simple-go-service
$ docker run --detach --network test-net --name go-3 simple-go-service
```

### Run Nginx Reverse Proxy
Run the nginx-reverse-proxy. Include `--network test-net` and also publish the listening port identified in the `nginx.conf` configuration file. Here port 8100 on the host machine is being mapped to port 8888 of the container. This will be the single point of entry to any of the simple-go-service containers.

```
$ docker run --detach --network test-net -p 8100:8888 --name proxy nginx-reverse-proxy
````

## Test Nginx Reverse Proxy
Check that all 4 containers are up and running.
```
$ docker ps
CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS                            NAMES
b31459ffab02        nginx-reverse-proxy   "nginx -g 'daemon of"    2 minutes ago       Up About a minute   80/tcp, 0.0.0.0:8100->8888/tcp   proxy
a124830c6bca        simple-go-service     "simple-go-service"      41 minutes ago      Up About a minute   8080/tcp                         go-3
1958a5a8b10c        simple-go-service     "simple-go-service"      41 minutes ago      Up About a minute   8080/tcp                         go-2
7fc356de9ad9        simple-go-service     "simple-go-service"      42 minutes ago      Up About a minute   8080/tcp                         go-1
```

Now try and connect to the proxy. `curl` is shown below but browser can be used instead. Calling `curl` or refreshing your browser should execute on different containers.

```

$ curl http://127.0.0.1:8100
Hello from container "7fc356de9ad9"
I have been visted 1 times
I was born "Tue Mar 13 23:15:39 UTC 2018"

$ curl http://127.0.0.1:8100
Hello from container "1958a5a8b10c"
I have been visted 1 times
I was born "Tue Mar 13 23:15:43 UTC 2018"

$ curl http://127.0.0.1:8100
Hello from container "a124830c6bca"
I have been visted 1 times
I was born "Tue Mar 13 23:15:50 UTC 2018"

```

Notice that the last 3 Go containers were started without `--rm`, as was the nginx one. This means that the containers are not removed after stopping. Instead of having to issue another `docker run` command, a `docker start <container name>` command can be use.

Try starting and stopping the Go containers. 

To delete a container issue `docker rm <container name>`

# Kubernetes
 Coming soon




