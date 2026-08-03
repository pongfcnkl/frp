.PHONY: dist install build preview lint

install:
	@cd .. && npm ci

build: install
	@cd .. && npm run build --workspace frps

dev:
	@npm run dev

preview:
	@npm run preview

lint:
	@npm run lint
