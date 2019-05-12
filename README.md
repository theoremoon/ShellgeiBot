# ShellGeiBot

## Build Docker image

```sh
$ ./build.bash shellgeibot:latest
```

## Test Docker image

```sh
$ docker container run --rm \
  -v $(pwd):/root/src \
  shellgeibot:latest \
  /bin/bash -c "apt update && apt install -y bats && bats /root/src/docker_image.bats"
```
