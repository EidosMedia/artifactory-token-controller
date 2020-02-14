FROM golang as build

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=build /go/src/app/artifactory-token-controller /

ENTRYPOINT ["/artifactory-token-controller"]