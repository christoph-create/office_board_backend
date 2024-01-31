#!/bin/sh
sudo docker image rm -f office-board-backend:latest
sudo docker build . -t office-board-backend
sudo docker run -d -p 8123:8123 -p 8124:8124 office-board-backend