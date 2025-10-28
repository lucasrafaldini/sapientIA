# Status do Roadmap v0.1

**Data de atualização**: 2025-10-27, 21:15  
**Milestone**: v0.1 (Alpha)  
**Prazo**: 2025-11-30

## Progresso Geral

```
v0.1: █████████░░░░░░ 43% (3/7 issues concluídas)
```

**Issues finalizadas**: 3 (#1, #2, #3)  
**Issues quase completas**: 1 (#6 @ 80%)  
**Issues pendentes**: 5

## Issues Criadas (7/7)

### ✅ Infraestrutura Básica - PARCIALMENTE COMPLETO

#### #6: Docker Compose (dev) + Make (DX) - 🟡 EM PROGRESSO

- **Status**: OPEN
- **Labels**: area:infra, kind:feature
- **Progresso atual**:
  - ✅ Makefile criado com alvos: build, dev, test, clean, install
  - ✅ docker-compose.yml estruturado (api, worker, ui, redis)
  - ⏳ Dockerfiles ainda não criados em `deployments/`
  - ⏳ Testes de ambiente completo pendentes

#### #7: Smoke test de pipeline - ⚪ NÃO INICIADO

- **Status**: OPEN
- **Labels**: area:infra, kind:feature
- **Progresso**: 0%
- **Bloqueio**: Depende de #1 (pipeline runner)

---

### ✅ Core/CLI - CONCLUÍDO

#### #1: CLI 'sapientia run' + parser de pipeline YAML - ✅ CONCLUÍDO

- **Status**: COMPLETO
- **Labels**: area:core, kind:feature
- **Progresso atual**:
  - ✅ CLI Cobra inicializado com comando `run`
  - ✅ go.mod criado e dependências instaladas
  - ✅ Build funcionando: `./bin/sapientia --version`
  - ✅ Parser YAML implementado (`internal/pipeline/pipeline.go`)
  - ✅ Executor completo implementado (`internal/pipeline/executor.go`)
  - ✅ Validações de pipeline robustas (tipos, nomes, estrutura)
  - ✅ Integração CLI + executor em `cmd/sapientia/run.go`
  - ✅ Testes unitários completos (9 testes passando)
  - ✅ Criação automática de diretórios de output
  - ✅ Placeholders para todos os tipos de step (ingest, lexical, graph, tree, report)

**Estimativa de conclusão**: 100% completo  
**Concluído em**: 2025-10-21

---

### ✅ Análise Lexical - CONCLUÍDO

#### #2: Léxico básico (freq, n-grams, TF-IDF) - ✅ CONCLUÍDO

- **Status**: COMPLETO
- **Labels**: area:core, kind:feature
- **Progresso**: 100%
- **Concluído em**: 2025-10-24
- **Detalhes**:
  - Tokenizador pt-BR com stopwords e filtro numérico (evita cantos/versos)
  - Freq, n-grams configuráveis, TF-IDF por documento
  - Export JSON + CSV opcionais; params: min_freq, ngrams, tfidf, stopwords_file
  - Testes unitários + E2E cobrindo pipelines de exemplo

---

### ✅ Grafos - CONCLUÍDO

#### #3: Grafo de coocorrência + export GEXF/CSV - ✅ CONCLUÍDO

- **Status**: COMPLETO
- **Labels**: area:core, kind:feature
- **Progresso**: 100%
- **Concluído em**: 2025-10-27
- **Detalhes**:
  - Coocorrência com janela configurável e peso mínimo de aresta
  - GEXF com pesos; CSVs de nós/arestas opcionais
  - Integração com pipeline, usando corpus/tokenização compartilhada
  - E2E em corpora grandes (Os Lusíadas, Divina Comédia)

---

### ⚪ Árvore Semântica - NÃO INICIADO

#### #4: Árvore semântica inicial - ⚪ NÃO INICIADO

- **Status**: OPEN
- **Labels**: area:core, kind:feature
- **Progresso**: 0%
- **Próximos passos**: Implementar `internal/tree/`

---

### ⚪ Relatórios - NÃO INICIADO

#### #5: Relatório HTML mínimo - ⚪ NÃO INICIADO

- **Status**: OPEN
- **Labels**: area:report, kind:feature
- **Progresso**: 0%
- **Próximos passos**: Implementar `internal/report/`

---

## Progresso Geral

```
Infra/DX:    [████████░░] 80% (estrutura criada, falta Dockerfiles)
CLI/Parser:  [███░░░░░░░] 30% (CLI funcional, parser pendente)
Léxico:      [░░░░░░░░░░]  0%
Grafos:      [░░░░░░░░░░]  0%
Árvore:      [░░░░░░░░░░]  0%
Relatório:   [░░░░░░░░░░]  0%
---
TOTAL v0.1:  [██░░░░░░░░] 18%
```

## Conquistas de Hoje (2025-10-21)

✅ **Estrutura completa de pastas** criada e segregada (Go/Python/Frontend)  
✅ **Go module** inicializado com CLI funcional  
✅ **Makefile** expandido com alvos úteis  
✅ **Docker Compose** estruturado  
✅ **Python worker** estruturado (FastAPI stub + requirements.txt)  
✅ **Frontend Next.js** estruturado (placeholder para v0.2)  
✅ **Pipeline de exemplo** criado (briefing.yaml)  
✅ **Dados de exemplo** adicionados  
✅ **.gitignore** atualizado para todos os ambientes  
✅ **Documentação** (STRUCTURE.md, READMEs por componente)

## Próximos Passos Críticos (Ordem de Prioridade)

1. **Issue #1**: Implementar parser YAML e executor de pipeline

   - Finalizar `internal/pipeline/pipeline.go`
   - Criar executor de DAG
   - Integrar com comando `run`

2. **Issue #6**: Criar Dockerfiles

   - `deployments/Dockerfile.go`
   - `deployments/Dockerfile.python`
   - `deployments/Dockerfile.nextjs`

3. **Issue #2**: — CONCLUÍDO

4. **Issue #3**: — CONCLUÍDO

5. **Issue #7**: Criar smoke test de pipeline

## Riscos e Bloqueios

⚠️ **Risco**: Issues #2-#7 dependem da conclusão de #1 (pipeline runner)  
⚠️ **Risco**: Integração Go ↔ Python worker ainda não testada  
⚠️ **Risco**: Performance em datasets médios não validada

## Métricas

- **Issues abertas**: 7/7 (100%)
- **Issues concluídas**: 0/7 (0%)
- **Dias até deadline**: ~40 dias
- **Velocity estimada**: precisa acelerar após #1

---

**Última atualização**: Estrutura básica completa, pronto para implementação dos módulos core.
