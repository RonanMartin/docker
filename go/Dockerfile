FROM golang as exec

COPY main.go /go/src/app/

ENV GO111MODULE=auto
ENV CGO_ENABLED=0

WORKDIR /go/src/app/

RUN go build -o app.go .

FROM alpine

WORKDIR /appexec

COPY --from=exec /go/src/app /appexec
RUN chmod -R 755 /appexec
EXPOSE 8080
ENTRYPOINT ./app.go
