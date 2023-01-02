FROM golang:1.15.2 as build

ENV SERVER_FILE /opt/social/social-app
ENV TESTGEN_FILE /opt/social/testdatagen

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${SERVER_FILE} cmd/social/main.go
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${TESTGEN_FILE} cmd/test_data_init/main.go

FROM alpine:3.9

LABEL ORGANIZATION="DMITRYI USOV INCORPRATED"
LABEL SERVICE="social"
LABEL MAINTAINERS="dimazusov@yandex.ru"

ENV SERVER_FILE "/opt/social/social-app"
ENV TESTGEN_FILE "/opt/social/testdatagen"

COPY --from=build ${SERVER_FILE} ${SERVER_FILE}
COPY --from=build ${TESTGEN_FILE} ${TESTGEN_FILE}

ENV CONFIG_FILE /etc/social/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}
COPY ./web /opt/social/web

WORKDIR /opt/social

CMD ${SERVER_FILE} -config ${CONFIG_FILE}
