FROM golang:1.22-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/wedback
RUN go build -o /wedback

FROM gcr.io/distroless/base-debian12:nonroot as release

LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"

WORKDIR /app/

COPY --from=builder /wedback .

EXPOSE 1378

ENTRYPOINT ["./wedback"]

CMD ["serve"]
