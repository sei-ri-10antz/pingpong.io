FROM alpine:latest

RUN apk --no-cache add tzdata
ADD ./build/server /server
ADD ./apis/v1/apidocs.swagger.json /apis/v1/apidocs.swagger.json

ENTRYPOINT [ "/server" ]