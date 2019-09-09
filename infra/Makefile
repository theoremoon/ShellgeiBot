IMAGE_NAME := shellgeibot_daemon

build:
	docker build -t $(IMAGE_NAME) .

start:
	docker run --rm -v $$PWD/../:/opt/ShellgeiBot -it $(IMAGE_NAME) bash

