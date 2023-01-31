FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /gh-issues-notifier


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gh-issues-notifier /gh-issues-notifier

VOLUME [ "/patterns.yaml" ]
VOLUME [ "/issues.txt" ]

EXPOSE 8080

ENTRYPOINT ["/gh-issues-notifier"]
