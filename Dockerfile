FROM golang:1.13.0-alpine AS build
RUN apk add bash ca-certificates git gcc g++ libc-dev make

WORKDIR /go/src/github.com/svazist/go-project-template
ENV GO111MODULE=on

ARG VERSION
ARG BUILD
ARG DATE

COPY . .

RUN go mod download

RUN make build

FROM alpine

RUN apk add ca-certificates

COPY --from=build /go/src/github.com/svazist/go-project-template/server /usr/bin/server
COPY --from=build /go/src/github.com/svazist/go-project-template/config.yaml /etc/config.yaml

EXPOSE 80

CMD [ "/usr/bin/server", "server" , "--config","/etc/config.yaml" ]
