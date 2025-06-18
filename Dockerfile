FROM golang:1.22 as builder
WORKDIR /src
COPY . .
RUN go build -o /out/psbl-watch ./cmd/psbl-watch

FROM gcr.io/distroless/static
COPY --from=builder /out/psbl-watch /psbl-watch
ENTRYPOINT ["/psbl-watch"]