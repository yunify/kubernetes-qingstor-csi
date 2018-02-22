FROM golang:alpine
MAINTAINER martinfan@yunify.com
ENV GOPATH /go/
ADD * /go/src/github.com/yunify/kubernetes-qingstor-csi/
WORKDIR /go/src/github.com/yunify/kubernetes-qingstor-csi/
RUN ls . && go build -o /bin/kubernetes-qingstor-csi
CMD ["/bin/kubernetes-qingstor-csi"]