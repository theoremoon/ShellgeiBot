# ShellGeiBot

## Twitter

- [@minyoruminyon](https://twitter.com/minyoruminyon)

## Specification

- https://furutsuki.hatenablog.com/entry/2018/07/13/221806

## Build Docker image

```sh
$ ./build.bash shellgeibot:latest
```

## Test Docker image

```sh
$ docker container run --rm \
  -v $(pwd):/root/src \
  shellgeibot:latest \
  /bin/bash -c "bats /root/src/docker_image.bats"
```
