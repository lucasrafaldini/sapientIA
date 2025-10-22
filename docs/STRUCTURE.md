# Estrutura do Projeto

```
sapientIA/
├── cmd/
│   └── sapientia/          # CLI principal (Cobra)
│       ├── main.go
│       └── run.go
│
├── internal/               # Código privado (não exportável)
│   ├── api/               # HTTP API (Chi) - v0.2
│   ├── pipeline/          # Parser e executor de pipelines YAML
│   ├── ingest/            # Ingestão e normalização de dados
│   ├── lexical/           # Análise lexical (freq, TF-IDF, n-grams)
│   ├── graph/             # Grafos de coocorrência
│   ├── tree/              # Árvores semânticas
│   ├── report/            # Geração de relatórios
│   └── worker/            # Cliente para Python worker
│
├── pkg/                   # Código público (exportável)
│   ├── models/           # Tipos e estruturas compartilhadas
│   └── utils/            # Utilitários gerais
│
├── worker/                # Python NLP Worker
│   ├── api/              # FastAPI endpoints
│   ├── nlp/              # Módulos NLP (spaCy, KeyBERT)
│   ├── models/           # Schemas Pydantic
│   └── requirements.txt
│
├── ui/                    # Next.js Frontend (v0.2)
│   ├── app/              # App Router pages
│   ├── components/       # Componentes React
│   ├── lib/              # Utilitários
│   └── package.json
│
├── data/                  # Dados de entrada
│   ├── examples/         # Exemplos para testes
│   ├── briefings/        # Gravações/textos de briefing
│   └── interviews/       # Respostas de entrevistas
│
├── pipelines/             # Definições de pipeline YAML
│   └── examples/
│       └── briefing.yaml
│
├── runs/                  # Outputs de execuções
│   └── .gitkeep
│
├── deployments/           # Dockerfiles
│   ├── Dockerfile.go
│   ├── Dockerfile.python
│   └── Dockerfile.nextjs
│
├── scripts/               # Scripts de automação
│
├── docs/                  # Documentação
│
├── .github/               # GitHub config (workflows, templates)
│
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Segregação por ambiente

### Go (Core)

- **Localização**: `cmd/`, `internal/`, `pkg/`
- **Responsabilidade**: Orquestração, pipelines, cálculos, persistência
- **Build**: `make build` → `bin/sapientia`
- **Teste**: `go test ./...`

### Python (Worker NLP)

- **Localização**: `worker/`
- **Responsabilidade**: NLP (lematização, embeddings, transcrição)
- **Setup**: `cd worker && pip install -r requirements.txt`
- **Execução**: `make dev-worker` ou `uvicorn api.main:app --reload`
- **Teste**: `pytest`

### Frontend (UI)

- **Localização**: `ui/`
- **Responsabilidade**: Interface web, visualização de grafos
- **Setup**: `cd ui && npm install`
- **Execução**: `cd ui && npm run dev`
- **Disponibilidade**: v0.2 (planejado)

## Comunicação entre ambientes

- **Go ↔ Python**: HTTP/gRPC (worker expõe API em `:8001`)
- **Go ↔ Frontend**: REST API em `:8080` (v0.2)
- **Shared data**: Volumes Docker (`data/`, `runs/`)

## Comandos úteis

```bash
# Build completo
make build

# Dev local (Docker)
make dev

# Testes
make test

# Clean
make clean

# Instalar CLI
make install
```
