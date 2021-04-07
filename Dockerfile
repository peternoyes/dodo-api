FROM golang:1.15

RUN mkdir -p /go/src
ADD . /go/src
WORKDIR /go/src

RUN go build -o dodo-api

CMD ["./dodo-api"]

EXPOSE 3000

RUN git clone https://github.com/cc65/cc65 /home/cc65

RUN cd /home/cc65 \
	&& make