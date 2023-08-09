FROM golang:1.20 as build
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /go/release
ADD . .
RUN CGO_ENABLED=0 go build .


FROM alpine  as prod
COPY --from=build /go/release/clash_airplan_filter  .
RUN chmod +x clash_airplan_filter
EXPOSE 8080
CMD ["./clash_airplan_filter"]