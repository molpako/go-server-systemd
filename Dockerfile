FROM ubuntu:latest
RUN apt update -y && apt install -y golang git init
CMD ["/sbin/init"]
