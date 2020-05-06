FROM golang:1.14.2-alpine3.11  as builder
RUN apk update && apk add git

WORKDIR /go/src/github.com/Kubernetes/ingress-nginx/images/custom-error-pages/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o /defaultbackend /go/src/github.com/Kubernetes/ingress-nginx/images/custom-error-pages/

FROM alpine:3.11
RUN apk update \
     && apk add ca-certificates \
     && rm -rf /var/cache/apk/* \
     && update-ca-certificates

EXPOSE 8080
ENTRYPOINT ["/defaultbackend", "--logtostderr=true"]
COPY --from=builder /defaultbackend /
COPY --from=builder /go/src/github.com/Kubernetes/ingress-nginx/images/custom-error-pages/roofts /roofts

