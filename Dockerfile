FROM golang:1.21

WORKDIR /COA

COPY . /COA

RUN apt update
# RUN apt install -y netcat-openbsd
# RUN apt install -y sockperf

RUN  go mod download \
    && go build /COA/main.go

ARG server_type=undefined

ENV SERVER_TYPE ${server_type}
ENV APP_ENVIRONMENT prod
ENV GROUP_UDP_ADDR 225.0.0.1:9000

# CMD  echo --server-type=$SERVER_TYPE --env=$APP_ENVIRONMENT GROUP_UDP_ADDR=$GROUP_UDP_ADDR \
CMD ./main --server-type $SERVER_TYPE --env $APP_ENVIRONMENT
