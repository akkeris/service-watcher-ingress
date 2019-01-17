FROM golang:1.10

RUN mkdir -p /go/src/service-watcher-ingress
ADD process.go  /go/src/service-watcher-ingress/process.go
ADD structs /go/src/service-watcher-ingress/structs
ADD services /go/src/service-watcher-ingress/services
ADD k8sconfig /go/src/service-watcher-ingress/k8sconfig
ADD utils /go/src/service-watcher-ingress/utils

ADD build.sh /build.sh
RUN chmod +x /build.sh
RUN /build.sh

RUN mkdir -p /root/.kube/certs
ADD start.sh /start.sh
RUN chmod +x /start.sh
CMD "/start.sh"



