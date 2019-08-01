# build
FROM golang:1.12.7-alpine3.10 as build

ENV PORT 8080
EXPOSE 8080

RUN mkdir /app
ADD . /app

ENV GOPROXY https://goproxy.io
ENV GIN_MODE release

WORKDIR  /app
RUN go build -tags=jsoniter -o user-test .


# release
FROM alpine:3.10
RUN mkdir /app
COPY --from=build /app/user-test /app/user-test

WORKDIR  /app
CMD ["/app/user-test"]
