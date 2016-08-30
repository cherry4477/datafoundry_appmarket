FROM golang:1.6.0

# for gateway
ENV SERVICE_NAME=datahub_stars

EXPOSE 8899

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

#RUN go get github.com/asiainfoLDP/datahub_stars
COPY . /go/src/github.com/asiainfoLDP/datahub_stars

WORKDIR /go/src/github.com/asiainfoLDP/datahub_stars

#RUN go get ./... && go build 

#RUN go get github.com/mattn/gom \
#    && gom install \
#    && gom build

#RUN go get github.com/tools/godep \
#    && $GOPATH/bin/godep restore \
#    && go build

#RUN go get github.com/tools/godep \
#    && godep go build   

RUN go build

CMD ["sh", "-c", "./datahub_stars -port=8899"]
