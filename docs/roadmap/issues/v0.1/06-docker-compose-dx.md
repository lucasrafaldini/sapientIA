# v0.1 — Docker Compose (dev) e DX

Labels sugeridas: `area:infra`, `kind:feature`

## Contexto

Ambiente reprodutível para desenvolvimento/testes locais.

## Objetivo

Configurar Docker Compose de desenvolvimento com targets Make convenientes.

## Escopo

- docker-compose.yml (serviços mínimos)
- Makefile com `make dev` e `make test`

## Critérios de aceite

- Subida local com `make dev` ou `docker compose up`

## Notas

- Ver `docs/roadmap/v0.1-execution.md` (Trilha 4)
