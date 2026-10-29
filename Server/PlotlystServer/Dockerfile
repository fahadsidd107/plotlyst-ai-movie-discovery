FROM golang:1.24-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/plotlyst-api .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/plotlyst-api /plotlyst-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/plotlyst-api"]
