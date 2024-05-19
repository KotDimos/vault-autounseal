FROM golang:1.21 AS build

WORKDIR /build

COPY ./go.mod ./go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/vault-autounseal


FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /build/bin/vault-autounseal /vault-autounseal

USER 10001
ENTRYPOINT ["/vault-autounseal"]
