FROM alpine:latest
RUN apk add --update alpine-sdk  bash sqlite sqlite-dev build-base gcc
RUN apk add --virtual .build-deps
RUN apk cache clean
ARG GO_LINK
RUN wget -P /root "$GO_LINK"
ARG GO_TAR
RUN tar -C /usr/local -xzf /root/"$GO_TAR"
ENV PATH=$PATH:/usr/local/go/bin
ENV CGO_ENABLED="1"
ENV APP_IP_ADDRESS="0.0.0.0"

# git config --global --add safe.directory [your dir here]
ENTRYPOINT ["bash"]
