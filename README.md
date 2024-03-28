# ImHost

[ImHost](./imhost) is a **web server that simply responds the "host name" of the server**.

Once the server is running, it listens to the port 80 and when accessed via HTTP, it returns the host name of the server.

Its main purpose is to be used for testing and debugging purposes in a Docker orchestration environment. Such as docker-compose, docker-swarm, kubernetes, etc.

We provide a [Dockerfile](./Dockerfile) for convenience. See the "[Docker Compose with scale and round-robin](#docker-compose-with-scale-and-round-robin)" section for the **full example using "--scale" option and Nginx as a load balancer** as well.

## Usage

### Docker

```shellsession
$ # Build the image
$ docker build -t imhost https://github.com/KEINOS/imhost.git
**snip**

$ # Run the container (expose the port 80 to 8080)
$ docker run --rm -p 8080:80 imhost
**snip**

$ # Check the hostname
$ curl -sS http://localhost:8080/
Hello from host: f9536b7afae0
```

### Docker Compose with scale and round-robin

Here is an example of `docker-compose.yml` that can:

1. Scale-up `imhost` containers using `--scale` option.
2. Load balancing using `nginx`.

```shellsession
$ # Run the containers and scale the imhost to 5 instances
docker compose up --scale imhost=5 --detach
**snip**

$ # Check the host names.
$ # Note that the host names are different each time because of
$ # the round-robin.
$ curl -sS http://localhost:8080
Hello from host: 2c80ea38d91a

$ curl -sS http://localhost:8080
Hello from host: 5fec68e38b9b

$ curl -sS http://localhost:8080
Hello from host: 353184aa84c0
```

```yaml
version: '3'

services:
  imhost:
    build:
      context: https://github.com/KEINOS/imhost.git
    expose:
      - "80"
    restart: unless-stopped

  loadbalancer:
    image: nginx:latest
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "8080:80"
    depends_on:
      - imhost
    restart: unless-stopped
```

```nginx
user nginx;

worker_processes auto;

events {
    worker_connections 1024;
}

http {
    server {
        listen 80;
        location / {
            proxy_pass http://imhost:80/;
        }
    }
}
```

- For the full example, see the [example](./_example) directory.

## Source code

- [/imhost/main.go](./imhost/main.go)

### Status

[![Unit Tests](https://github.com/KEINOS/imhost/actions/workflows/unit-test.yml/badge.svg)](https://github.com/KEINOS/imhost/actions/workflows/unit-test.yml)
[![GolangCI-Lint Test](https://github.com/KEINOS/imhost/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/KEINOS/imhost/actions/workflows/golang-ci.yml)
[![Docker Tests](https://github.com/KEINOS/imhost/actions/workflows/docker-test.yml/badge.svg)](https://github.com/KEINOS/imhost/actions/workflows/docker-test.yml)
[![CodeQL](https://github.com/KEINOS/imhost/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/KEINOS/imhost/actions/workflows/github-code-scanning/codeql)
[![codecov](https://codecov.io/gh/KEINOS/imhost/graph/badge.svg?token=7WsjthYoE6)](https://codecov.io/gh/KEINOS/imhost)