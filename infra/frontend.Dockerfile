FROM --platform=$BUILDPLATFORM node:24-alpine AS deps
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

FROM deps AS build
COPY frontend/ ./
ARG DISCORD_CLIENT_ID
RUN DISCORD_CLIENT_ID="$DISCORD_CLIENT_ID" npm run build

FROM nginxinc/nginx-unprivileged:stable-alpine AS runtime
COPY infra/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
