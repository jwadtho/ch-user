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
ADD fluent-bit.conf .
ADD parsers.conf .
ADD output-datadog.conf .
RUN mkdir "dd-output-config"

ENTRYPOINT ["k8s-tags"]