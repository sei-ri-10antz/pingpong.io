FROM alpine:latest

RUN apk --no-cache add tzdata
ADD ./build/server /server

ENTRYPOINT [ "/server" ]