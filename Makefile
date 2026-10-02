dev:
	$(MAKE) -j2 dev-frontend dev-backend

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run ./cmd/api

build:
	cd frontend && npm run build
	cd backend && go build -o ../bin/server ./cmd/api

test:
	cd backend && go test ./...
	cd frontend && npm run typecheck

.PHONY: dev dev-frontend dev-backend build test
