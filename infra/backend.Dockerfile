FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS deps
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download

FROM deps AS build
ARG TARGETOS TARGETARCH
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /server ./cmd/api

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 app
COPY --from=build /server /server
USER app
ENTRYPOINT ["/server"]
