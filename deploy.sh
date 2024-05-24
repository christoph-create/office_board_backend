#!/bin/sh
sudo docker image rm -f office-board-backend:latest
sudo docker build . -t office-board-backend
sudo docker run -it --rm --mount type=bind,src="/home/golo",target=/mnt -p 8123:8123 office-board-backend