FROM ubuntu:latest

RUN apt update -y && apt install -y init git wget vim
RUN wget https://go.dev/dl/go1.19.linux-arm64.tar.gz && tar -xvf go1.19.linux-arm64.tar.gz && mv go /usr/local

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

ADD ./systemd/hello.service /etc/systemd/system/hello.service
ADD ./systemd/hello.socket /etc/systemd/system/hello.socket

WORKDIR /app
ADD ./main.go ./go.mod ./go.sum /app/
RUN go build -o go-server-systemd ./main.go
RUN systemctl enable hello.socket

CMD ["/sbin/init"]
