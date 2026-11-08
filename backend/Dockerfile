FROM node:26-alpine AS frontend
WORKDIR /frontend
COPY frontend/xpense-trackr-frontend/package.json frontend/xpense-trackr-frontend/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY frontend/xpense-trackr-frontend/ ./
ARG REACT_APP_API_URL=
ENV REACT_APP_API_URL=${REACT_APP_API_URL}
RUN npm run build

FROM golang:1.27-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/xpense-trackr ./cmd/expense-tracker

FROM alpine:3.22
RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /data /app/web \
    && chown -R app:app /data /app
WORKDIR /app
COPY --from=backend /out/xpense-trackr /app/xpense-trackr
COPY --from=frontend --chown=app:app /frontend/build/ /app/web/
ENV APP_ADDR=:8080 \
    DB_DRIVER=sqlite \
    DB_PATH=/data/xpense-trackr.db \
    APP_DATA_DIR=/data \
    WEB_DIR=/app/web
VOLUME ["/data"]
EXPOSE 8080
USER app
ENTRYPOINT ["/app/xpense-trackr"]
