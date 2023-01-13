# syntax=docker/dockerfile:1.3-labs
FROM golang:1.19-alpine as builder

COPY ./ /src
WORKDIR /src
RUN  CGO_ENABLED=0 go build -o bin/kifthenfi -ldflags "-s -w" .

FROM alpine
COPY --from=builder /src/bin/kifthenfi /usr/local/bin/kifthenfi
ENTRYPOINT [ "/usr/local/bin/kifthenfi" ]



