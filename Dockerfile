FROM golang as builder
WORKDIR $GOPATH/src/github.com/metalstormbass/mike-admission-controller
COPY ./ .
RUN go build -ldflags="-w -s" -v
RUN cp ca.crt / && cp ca.key /
RUN cp mike-admission-controller /


FROM ubuntu
COPY --from=builder mike-admission-controller .
COPY --from=builder ca.crt .
COPY --from=builder ca.key .
CMD ./mike-admission-controller 