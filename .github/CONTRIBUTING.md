# Contribuindo para o SapientIA

Obrigado por contribuir! Este documento descreve como participar do projeto.

## Requisitos
- Go 1.22+
- Python 3.10+ (para o worker NLP)
- Node 20+ (UI)
- Docker e Docker Compose
- Make / GNU tools

## Como começar
1. Faça um fork e clone o repositório.
2. Crie uma branch: `git checkout -b feat/minha-feature`.
3. Suba o ambiente local:
   ```bash
   make dev
   # ou
   docker compose up
   ```
4. Rode os testes:
   ```bash
   make test
   ```

## Padrão de commits
Use o formato do Conventional Commits:

```
type(scope): descrição curta
```

Exemplos:

```text
feat(pipeline): add YAML parser
fix(graph): correct PMI calculation
docs(readme): update usage section
```

## Código e estilo
- Go: `gofmt`, `golangci-lint`
- Python: `ruff`, `pytest`
- JS/TS: `eslint`, `prettier`

## Testes
- Toda função pública deve possuir ao menos um teste.
- PRs sem testes podem receber solicitação para adicioná-los.

## Discussões e issues
- Use Issues para bugs e propostas de feature.
- Descreva o contexto, cenário de uso e passos para reproduzir.

## Licenças
- Core do projeto: Apache-2.0.
- Partes Enterprise: EULA (ver `/licenses`).
- Não submeter código incompatível com Apache-2.0 para o Core.

Obrigado por ajudar a construir uma ferramenta que escuta para transformar. 💚