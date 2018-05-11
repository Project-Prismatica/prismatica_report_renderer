FROM golang:alpine AS build-env
WORKDIR /usr/local/go/src/github.com/Project-Prismatica/prismatica_report_renderer
COPY . /usr/local/go/src/github.com/Project-Prismatica/prismatica_report_renderer
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go get -u ./...
RUN go build -o build/prismatica_report_renderer ./prismatica_report_renderer


FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build-env /usr/local/go/src/github.com/Project-Prismatica/prismatica_report_renderer/build/prismatica_report_renderer /bin/prismatica_report_renderer
CMD ["prismatica_report_renderer", "up"]
