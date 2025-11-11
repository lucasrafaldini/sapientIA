# 🌿 SapientIA

[![Build](https://img.shields.io/github/actions/workflow/status/yourorg/sapientia/ci.yml?branch=main&label=build)](https://github.com/yourorg/sapientia/actions)
[![License](https://img.shields.io/github/license/yourorg/sapientia)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yourorg/sapientia)](go.mod)
[![Contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)

> “_Sapientia est intellegentia quae audit._”  
> **Sabedoria é a inteligência que escuta.**

---

## O que é o SapientIA

**SapientIA** é uma plataforma open-source para **análise de conversas e dados textuais corporativos**, desenhada para substituir o fluxo fragmentado com **Excel + Gephi + IRaMuTeQ**.  
Transforma briefs, entrevistas e respostas abertas em **árvores semânticas**, **redes de coocorrência** e **insights acionáveis** — uma verdadeira **terapia de dados** para organizações.

**Propósito:** dar às empresas um espaço de escuta. Onde havia ruído e dispersão, surge clareza, empatia e estratégia.

---

## Como funciona

1. **Briefing**  
   Faça upload de uma gravação de reunião ou texto-base. O sistema transcreve, identifica dores e temas e gera uma **árvore semântica** do discurso.

2. **Perguntas-chave**  
   A partir da árvore do briefing, o SapientIA sugere ~**10 perguntas** que cobrem os ramos prioritários.

3. **Entrevistas**  
   Colete respostas de dezenas de participantes. Para **cada pergunta**, o sistema gera **uma árvore semântica** a partir das respostas agregadas.

4. **Análise**  
   Descubra padrões lexicais, termos característicos (χ²/PMI), relações entre temas e **insights por pergunta** e **transversais**.

5. **Relatório**  
   Gere **HTML/PDF** com tabelas, gráficos e evidências (trechos ancorados), além de exportar **CSV/JSON/GEXF/SVG**.

---

## Stack (v1)

| Camada         | Tecnologia                              | Função                                    |
| -------------- | --------------------------------------- | ----------------------------------------- |
| **Core**       | Go 1.22+, Gonum, SQLC                   | Pipelines, grafos, cálculo e persistência |
| **Worker**     | Python (spaCy pt-BR, KeyBERT, Whisper)  | Lematização, embeddings e transcrição     |
| **API**        | Chi (Go) + Asynq (Redis)                | Orquestração de jobs e serviços           |
| **Frontend**   | Next.js, Tailwind, Cytoscape.js         | Dashboards e grafos interativos           |
| **Relatórios** | Go templates + (wkhtmltopdf/WeasyPrint) | Exportações HTML/PDF                      |
| **Infra**      | Docker Compose                          | Ambiente reprodutível                     |

> **Visão de longo prazo:** manter **binário único** (Go) que auto-extrai o _worker_ e baixa modelos sob demanda.  
> **IA avançada (LLMs) prevista** para v5/v6, após validação com clientes.

---

## Começando

```bash
# Clone o repositório
git clone https://github.com/yourorg/sapientia.git
cd sapientia

# Instale o binário principal (CLI)
go install ./cmd/sapientia

# Execute um pipeline de exemplo
sapientia run pipelines/exemplo.yaml

# Escolha um perfil de stopwords (core/extended/aggressive)
sapientia run pipelines/exemplo.yaml --stopwords extended
```

## Estrutura do repositório

```bash
sapientia/
├─ cmd/                # CLI (Cobra)
├─ internal/
│  ├─ api/             # HTTP (Chi)
│  ├─ pipeline/        # Parser YAML + DAG
│  ├─ ingest/          # Transcrição, limpeza, normalização
│  ├─ lexical/         # Freq, tf-idf, χ², n-grams, AFC/SVD
│  ├─ graph/           # Coocorrência, comunidades, export GEXF/GraphML
│  ├─ tree/            # Hierarquização para "árvore semântica"
│  ├─ report/          # Templates e exportações (HTML/PDF/SVG)
│  └─ worker/          # Integração com worker Python (gRPC/HTTP)
├─ ui/                 # Next.js + Tailwind + Cytoscape.js
├─ pipelines/          # Pipelines de exemplo (YAML)
└─ docs/               # Documentação (MkDocs)
```

## Roadmap

| Versão | Foco principal                                       | Status          | Progresso        |
| -----: | ---------------------------------------------------- | --------------- | ---------------- |
|   v0.1 | CLI + pipeline YAML + coocorrência + export GEXF     | 🔄 Em dev       | 43% (3/7 issues) |
|   v0.2 | API Go + UI web + grafos interativos                 | ⏳ Planejado    | -                |
|   v0.3 | χ², AFC/SVD, relatórios completos (HTML/PDF)         | ⏳ Planejado    | -                |
|   v1.0 | Binário único (Go + worker embutido + model manager) | 🎯 Meta release | -                |
|  v5/v6 | LLMs para insights/assistência analítica             | 🔮 Longo prazo  | -                |

### v0.1 - Issues Concluídas

- ✅ **#1**: CLI 'sapientia run' + parser de pipeline YAML + executor ([PR #8](https://github.com/lucasrafaldini/sapientIA/pull/8))
- ✅ **#2**: Léxico básico (freq, n-grams, TF-IDF) com stopwords e CSV ([PR #9](https://github.com/lucasrafaldini/sapientIA/pull/9))
- ✅ **#3**: Grafo de coocorrência + export GEXF/CSV com stopwords configuráveis ([PR #10](https://github.com/lucasrafaldini/sapientIA/pull/10))

➡️ Veja o roadmap completo e critérios por versão em [ROADMAP.md](ROADMAP.md).  
📋 Acompanhe o progresso da v0.1 no [Project Board](https://github.com/users/lucasrafaldini/projects/1).  
📊 Status detalhado em [docs/roadmap/STATUS.md](docs/roadmap/STATUS.md).

## Exemplos de uso (CLI)

```bash
# 1) Rodar pipeline completo (briefing → perguntas → entrevistas → árvores → relatório)
sapientia run pipelines/cliente-x.yaml

# 2) Rodar com perfil de stopwords específico (core, extended, aggressive)
sapientia run pipelines/cliente-x.yaml --stopwords extended

# 3) Gerar só a árvore do briefing
sapientia briefing build --audio data/briefing/gravacao.mp3 --out runs/cliente-x/

# 4) Sugerir ~10 perguntas com cobertura de ramos
sapientia questions make --tree runs/cliente-x/briefing_tree.json --target 10

# 5) Ingerir entrevistas em lote
sapientia interviews ingest --csv data/entrevistas/respostas.csv

# 6) Gerar árvore por pergunta
sapientia question-tree build --question-id Q3 --out runs/cliente-x/Q3/

# 7) Montar relatório final
sapientia report build --workspace runs/cliente-x/ --format html,pdf
```

### Stopwords

O SapientIA inclui três perfis de stopwords pt-BR:

- **core** (padrão): artigos, pronomes, preposições, conjunções
- **extended**: core + advérbios frequentes e auxiliares (ser, estar, ter, haver)
- **aggressive**: extended + verbos de alta frequência (fazer, dizer, ir, poder, etc.)

Use `--stopwords` para sobrescrever o perfil em qualquer pipeline. Veja [data/stopwords/README.md](data/stopwords/README.md) para detalhes.

## ⚖️ Licença

O núcleo do **SapientIA** é distribuído sob a licença **Apache-2.0**.  
Recursos avançados (ex: módulos de NLP, UI web, relatórios corporativos)  
podem ser licenciados separadamente sob a **SapientIA Enterprise License**.

## Como contribuir

Contribuições são bem-vindas! Siga estes passos:

Faça um fork do repositório

Crie uma branch: git checkout -b feat/minha-feature

Faça commits claros: git commit -m "feat: adiciona X"

Abra um Pull Request 🧡

Veja também: CONTRIBUTING.md (padrões de código, testes, ambiente, commits).

## Comunidade

Issues & Roadmap: https://github.com/yourorg/sapientia/issues
