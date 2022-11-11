FROM golang:1.19.3-alpine3.15 as builder
COPY go.mod go.sum /go/src/gitlab.com/shreyanshyad/task/
WORKDIR /go/src/gitlab.com/shreyanshyad/task
RUN go mod download
COPY . /go/src/gitlab.com/shreyanshyad/task
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/run_server /go/src/gitlab.com/shreyanshyad/task/dashboard

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/shreyanshyad/task/build/run_server /usr/bin/run_server
COPY ./dashboard/db/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/usr/bin/run_server"]