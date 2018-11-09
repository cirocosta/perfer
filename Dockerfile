FROM alpine

RUN set -x && \
	apk add \
		--update-cache \
		--repository http://dl-3.alpinelinux.org/alpine/edge/testing \
		git bash perf perl

RUN set -x && \
	git clone https://github.com/brendangregg/FlameGraph

ADD ./tools /usr/local/bin

ENV PATH=/FlameGraph:$PATH

