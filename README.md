# ShellGeiBot

## Build Docker image

```sh
$ ./build.bash
```

## Test Docker image

- Uses [Bats](https://github.com/sstephenson/bats) for test docker image.

### Installation

- Linux (with APT)

```sh
$ sudo apt install bats
```

- macOS (with Homebrew)

```sh
$ brew install bats
```

### Run

```sh
$ bats docker_image.bats
```
