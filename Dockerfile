FROM golang:1.22 AS build

WORKDIR /build

COPY ./go.mod ./go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/vault-autounseal


FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /build/bin/vault-autounseal /vault-autounseal

USER 1001
ENTRYPOINT ["/vault-autounseal"]
