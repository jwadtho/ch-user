FROM golang:1.17-alpine as dev-env

WORKDIR /app

FROM dev-env as build-env
COPY . /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -o /k8s-tags

FROM alpine:3.14.1 as runtime

COPY --from=build-env /k8s-tags /usr/local/bin/k8s-tags
RUN chmod +x /usr/local/bin/k8s-tags

ENTRYPOINT ["k8s-tags"]