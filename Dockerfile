FROM golang
ADD . /go/src/github.com/vsco/dcdr
WORKDIR /go/src/github.com/vsco/dcdr
RUN ./script/install
ADD ./config/config.example.hcl /etc/dcdr/config.hcl
ENTRYPOINT ["/go/bin/dcdr"]

