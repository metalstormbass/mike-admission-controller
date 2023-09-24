FROM cgr.dev/chainguard/go AS builder
WORKDIR $GOPATH/src/github.com/metalstormbass/mike-admission-controller
COPY ./ .
RUN go build -ldflags="-w -s" -v
RUN cp server.crt / && cp server.key /
RUN cp mike-admission-controller /

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder mike-admission-controller .
COPY --from=builder server.crt .
COPY --from=builder server.key .
CMD ["/mike-admission-controller"] 