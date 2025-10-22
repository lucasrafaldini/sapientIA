# 🧩 Arquitetura SapientIA (detalhada)

## Backend (Go)
- CLI: Cobra / Viper (YAML)
- API: Chi (handlers REST)
- Jobs: Asynq (Redis)
- Dados: SQLite (MVP) → Postgres
- Grafos: Gonum (estruturas e métricas)
- Relatórios: Go templates + wkhtmltopdf

## NLP Worker (Python)
- spaCy pt_core_news_lg (lematização/PoS)
- KeyBERT/YAKE (keyphrases)
- Whisper/faster-whisper (transcrição)
- prince (AFC/ACM)
- Exposição via HTTP/gRPC

## UI (Next.js)
- Cytoscape.js (grafos)
- Recharts (gráficos)
- Tailwind + shadcn/ui (design)

## Árvore semântica (algoritmo)
1. Limpeza + lemas + n-grams
2. Grafo de coocorrência (janela/sentença)
3. Comunidades (Louvain/Leiden)
4. Similaridade entre comunidades
5. Clustering hierárquico (Ward/average) → árvore
6. Edição manual (mesclar/renomear ramos)

## Binário único
- Go embute UI estática e bootstrapper.
- Na 1ª execução: auto-extrai `worker-nlp` (Nuitka/PyInstaller) e baixa modelos em `~/.sapientia/models`.

## Pastas (resumo)
- `/cmd` CLI
- `/internal/pipeline` DAG e estágios
- `/internal/ingest` transcrição/normalização
- `/internal/lexical` contagens, χ², PMI, SVD/AFC
- `/internal/graph` coocorrência, comunidades, export
- `/internal/report` HTML/PDF
- `/internal/worker` client e protocolos
- `/ui` Next.js
