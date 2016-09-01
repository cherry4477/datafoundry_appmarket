FROM golang:1.6.2

EXPOSE 8000

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

COPY . /go/src/github.com/asiainfoLDP/datafoundry_appmarket

WORKDIR /go/src/github.com/asiainfoLDP/datafoundry_appmarket

RUN go build

CMD ["sh", "-c", "./datafoundry_appmarket -port=8000"]
