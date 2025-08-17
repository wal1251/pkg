ARG GOLANG_VERSION=1.22-alpine

FROM registry.redmadrobot.com:5005/backend-go/rmr-pkg/golang:$GOLANG_VERSION

RUN apk -U add ca-certificates
RUN apk update && apk add git
RUN go install -v golang.org/x/tools/cmd/godoc@latest

WORKDIR /go

COPY . src

ENTRYPOINT ["godoc", "-index", "-goroot=/go", "-http=:6060"]
