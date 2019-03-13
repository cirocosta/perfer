FROM alpine

RUN apk add \
	--update-cache \
	--repository http://dl-3.alpinelinux.org/alpine/edge/testing \
	git bash perf perl

RUN git clone https://github.com/brendangregg/FlameGraph

ADD ./tools /usr/local/bin

ENV PATH=/FlameGraph:$PATH

