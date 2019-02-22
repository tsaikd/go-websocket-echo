FROM golang as builder
ARG HTTP_PROXY="$HTTP_PROXY"
ARG HTTPS_PROXY="$HTTPS_PROXY"
ENV GO_WEBSOCKET_ECHO_HOME="/go_websocket_echo"
ENV HTTP_PROXY="$HTTP_PROXY"
ENV HTTPS_PROXY="$HTTPS_PROXY"
ENV CGO_ENABLED="0"
WORKDIR $GO_WEBSOCKET_ECHO_HOME
COPY "build/certs/*.crt" "/usr/local/share/ca-certificates/"
RUN rm -f "/usr/local/share/ca-certificates/empty.crt"
RUN update-ca-certificates
COPY go.mod $GO_WEBSOCKET_ECHO_HOME
COPY go.sum $GO_WEBSOCKET_ECHO_HOME
RUN go mod download
RUN gobuilder version -c ">=0.2.1" >/dev/null 2>&1 || go get -v "github.com/tsaikd/gobuilder"
COPY . $GO_WEBSOCKET_ECHO_HOME
RUN gobuilder

FROM alpine
EXPOSE 8080
COPY --from=builder "/go_websocket_echo/go-websocket-echo" "/usr/local/bin/websocket-echo"
RUN chmod +x "/usr/local/bin/websocket-echo"
CMD websocket-echo server
