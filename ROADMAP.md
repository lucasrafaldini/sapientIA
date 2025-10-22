# 🗺️ Roadmap do SapientIA

Este roadmap detalha os marcos planejados, critérios de conclusão e prioridades do projeto. Ele expande o resumo presente no `README.md` e serve como referência para criação de Issues e Milestones no GitHub.

## Princípios de priorização

- Valor para o usuário primeiro: ciclo completo de “briefing → perguntas → entrevistas → árvores → relatório”.
- Simplicidade e reprodutibilidade: execução local via Docker Compose e binários fáceis de instalar.
- Qualidade acima de quantidade: testes, métricas e documentação mínima por feature.
- Segurança e privacidade por padrão: dados locais por default; integrações opt-in.

## Visão geral por versão

| Versão | Foco principal                                       | Status         | Progresso    |
| -----: | ---------------------------------------------------- | -------------- | ------------ |
|   v0.1 | CLI + pipeline YAML + coocorrência + export GEXF     | 🔄 Em dev      | 31% (2.17/7) |
|   v0.2 | API Go + UI web (gráficos interativos) + jobs/queue  | ⏳ Planejado   | -            |
|   v0.3 | χ²/PMI/AFC/SVD + relatórios HTML/PDF completos       | ⏳ Planejado   | -            |
|   v0.4 | Persistência (Postgres), observabilidade, e2e        | ⏳ Planejado   | -            |
|   v0.5 | Ingestões/integrações (Drive/S3), STT Whisper, spaCy | ⏳ Planejado   | -            |
|   v0.6 | Performance, UX e RBAC                               | ⏳ Planejado   | -            |
|   v1.0 | Binário único (Go + worker embutido/auto-download)   | 🎯 Meta        | -            |
|  v5/v6 | LLMs para insights/assistência analítica             | 🔮 Longo prazo | -            |

---

## v0.1 — Alpha (Fundação do Core)

Objetivo: pipeline reprodutível fim-a-fim com exportações úteis para analistas.

Escopo mínimo:

- CLI (`cmd/sapientia`) com subcomandos essenciais e parser de pipeline YAML (`internal/pipeline`).
- Coocorrência e grafos básicos (`internal/graph`), export GEXF/GraphML.
- Módulos lexicais iniciais (`internal/lexical`): frequências, n-grams, TF-IDF simples.
- Árvore semântica inicial (`internal/tree`) a partir de temas/termos.
- Relatórios mínimos (HTML simples) e CSV/JSON de outputs.
- Docker Compose de desenvolvimento; Make targets (`make dev`, `make test`).
- Worker Python stub acessível (HTTP/gRPC) para preparar integração futura.

Critérios de conclusão:

- Rodar `sapientia run pipelines/exemplo.yaml` gera artefatos (GEXF/CSV/JSON/HTML) estáveis.
- Testes: ≥60% das funções públicas cobertas (unitários básicos em Go + smoke tests no pipeline).
- Documentação “Começando” no README e `CONTRIBUTING.md` saneado.

Riscos/mitigações:

- Complexidade de parsing YAML → validar com esquemas e exemplos.
- Performance em datasets médios → amostragens e estruturas de dados eficientes (Gonum, buffers).

👉 Detalhamento técnico e checklist: veja o plano de execução da v0.1 em [`docs/roadmap/v0.1-execution.md`](docs/roadmap/v0.1-execution.md).  
📊 Status e progresso atualizado: [`docs/roadmap/STATUS.md`](docs/roadmap/STATUS.md).

### ✅ Issues Concluídas na v0.1

- **#1**: CLI 'sapientia run' + parser de pipeline YAML + executor (concluída em 21/10/2025)

---

## v0.2 — MVP API + UI

Objetivo: orquestrar jobs e visualizar grafos no navegador.

Escopo:

- API HTTP em Go (`internal/api`) com filas de jobs (Asynq/Redis) e status de execução.
- UI Next.js (`ui/`) com páginas: upload/briefing, perguntas, entrevistas, análise e resultados.
- Gráficos com Cytoscape.js; exibição de métricas básicas nos nós/arestas.
- Autenticação simples (token/env ou sessão dev) para ambiente local.

Critérios:

- Usuário consegue subir um briefing, ver árvore/grafo, acompanhar progresso do job e baixar exportações.
- Testes de integração básicos (API + UI stubbed) e health checks.

---

## v0.3 — Análises avançadas e Relatórios

Objetivo: entregar valor analítico com profundidade e relatórios executivos.

Escopo:

- Estatísticas: χ², PMI calibrado, AFC/SVD para reduções/associações.
- Relatórios HTML/PDF com templates robustos (`internal/report`), evidências ancoradas e SVGs.
- Exportações: CSV/JSON/SVG adicionais e metadados para auditoria.

Critérios:

- Relatório final com sumário, seções por pergunta, destaques transversais e anexos técnicos.
- Testes de regressão nos cálculos estatísticos e snapshots de relatório.

---

## v0.4 — Persistência e Operabilidade

Objetivo: tornar o sistema operável e confiável para projetos maiores.

Escopo:

- Postgres para runs/artefatos; migrações versionadas.
- Observabilidade (logs estruturados, métricas, traces opcionais) e retries/timeout nos jobs.
- Testes e2e no Docker Compose; passos de backup/restore de dados.

Critérios:

- Reprocessamento idempotente e retentativa automática em falhas transitórias.
- Dashboards mínimos de saúde (UI ou endpoints) e docs de operação.

---

## v0.5 — Ingestões e Integrações

Objetivo: reduzir atrito de entrada de dados e ampliar casos de uso.

Escopo:

- Conectores: Google Drive/Sheets, S3; importadores CSV/JSON robustos.
- STT Whisper integrado no pipeline de briefing (Python worker).
- Gestão de modelos spaCy pt-BR e dicionários; cache de embeddings.

Critérios:

- Upload ou seleção de fonte externa dispara pipeline completo com monitoramento.
- Documentação das integrações e quotas.

---

## v0.6 — Performance, UX e RBAC

Objetivo: polir a experiência e suportar volumes maiores.

Escopo:

- Otimizações de performance (indexação, streaming de progresso, chunking).
- UX: estados de carregamento, filtros e buscas nos grafos/árvores.
- Controles de acesso (RBAC leve) para espaços de trabalho.

Critérios:

- Conjuntos de dados 10–50× maiores processando em tempos previsíveis.
- UI fluida com feedback de progresso e filtros responsivos.

---

## v1.0 — Binário único e distribuição

Objetivo: distribuição simples para times não técnicos.

Escopo:

- Binário Go que auto-extrai/baixa o worker e modelos sob demanda.
- Instaladores/embalagem e checagens de ambiente.
- Guias de migração e SLOs de release.

Critérios:

- Instalação em “um passo” e execução de um pipeline de exemplo sem fricção.

---

## Backlog (Próximos / Depois)

- Insights assistidos por LLM (resumos por pergunta, hipóteses, QA sobre dados).
- Modelagem de tópicos (BERTopic, LDA) comparativa com abordagem de árvores.
- Colaboração: comentários ancorados e marcações no relatório.
- Internacionalização completa (pt/en/es) e modelos multi-idioma.
- Plugins de exportação (PowerPoint/Keynote) e templates corporativos.

---

## Riscos e mitigação

- Licenças/modelos de NLP: padronizar em modelos permissivos e documentar custos.
- Privacidade: dados locais por default; anonimização opcional.
- Sustentação do stack: versões pinadas e CI com matrix para Go/Python/Node.

## Métricas de sucesso

- Time-to-Insight (do upload ao primeiro gráfico/relatório).
- Adoção: execuções/mês e número de usuários ativos.
- Qualidade: cobertura de testes, taxa de falhas por run, tempos p95.

## Como contribuir para o roadmap

1. Abra uma Issue descrevendo o contexto, a necessidade e o impacto.
2. Atribua labels sugeridas: `area:core`, `area:ui`, `area:worker`, `kind:feature`, `kind:bug`, `good first issue`.
3. Relacione a Issue à versão/marco (Milestone) correspondente.
4. Traga discussões para o Discord/Issues e mantenha a conversa aberta.

Sugestão de labels:

- `area:core`, `area:api`, `area:ui`, `area:worker`, `area:report`, `area:infra`
- `kind:feature`, `kind:bug`, `kind:docs`, `kind:perf`, `kind:refactor`
- `good first issue`, `help wanted`, `blocked`
