FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/cmd/server/main /app/
WORKDIR /app
ENTRYPOINT ["./main"]