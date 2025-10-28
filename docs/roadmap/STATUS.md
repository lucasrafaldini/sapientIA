# Status do Roadmap v0.1

**Data de atualização**: 2025-10-21, 21:15  
**Milestone**: v0.1 (Alpha)  
**Prazo**: 2025-11-30

## Progresso Geral

```
v0.1: ███████░░░░░░░░ 31% (2,17/7 issues concluídas)
```

**Issues finalizadas**: 1 (#1)  
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

### ⚪ Análise Lexical - NÃO INICIADO

#### #2: Léxico básico (freq, n-grams, TF-IDF) - ⚪ NÃO INICIADO

- **Status**: OPEN
- **Labels**: area:core, kind:feature
- **Progresso**: 0%
- **Próximos passos**: Implementar `internal/lexical/`

---

### ⚪ Grafos - NÃO INICIADO

#### #3: Grafo de coocorrência + export GEXF/GraphML/CSV/JSON - ⚪ NÃO INICIADO

- **Status**: OPEN
- **Labels**: area:core, kind:feature
- **Progresso**: 0%
- **Próximos passos**: Implementar `internal/graph/`

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

3. **Issue #2**: Implementar módulo lexical básico

   - Frequências
   - N-grams
   - TF-IDF

4. **Issue #3**: Implementar grafos de coocorrência

   - Matriz de coocorrência
   - Export GEXF/GraphML

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
