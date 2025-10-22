# 🧭 Roadmap — SapientIA

> Última atualização: 2025-10-21  
> Este arquivo é lido por agentes de desenvolvimento e serve como guia de direção técnica e estratégica.

---

## 🎯 Visão geral

O **SapientIA** é uma plataforma open-source para análise semântica de conversas corporativas, estruturada em estágios progressivos.  
Cada versão tem um _tema de foco_ e uma _meta técnica_.

---

## 🧱 Estrutura do roadmap

| Versão    | Foco principal                                | Estado                | Observações                                         |
| --------- | --------------------------------------------- | --------------------- | --------------------------------------------------- |
| **v0.1**  | CLI + pipelines YAML + análise lexical básica | 🚧 Em desenvolvimento | MVP funcional via terminal                          |
| **v0.2**  | API e UI web (Next.js) + grafos interativos   | ⏳ Planejado          | Interface web e visualização das árvores            |
| **v0.3**  | χ² / AFC / relatórios automáticos HTML-PDF    | ⏳ Planejado          | Relatórios interpretativos e exportação corporativa |
| **v1.0**  | Binário único (Go + worker embutido)          | ⏳ Planejado          | Entrega autônoma (sem dependências externas)        |
| **v5/v6** | IA generativa e insights automáticos          | 🌙 Longo prazo        | Modelos LLM fine-tunados para análise semântica     |

---

## 🚀 Sprint atual (v0.1 — MVP)

**Meta:** CLI funcional que processa um pipeline completo (`sapientia run`)  
**Prazo sugerido:** 4 semanas

### Subtarefas prioritárias

- [ ] Implementar `internal/pipeline` (parser YAML → DAG)
- [ ] Implementar `internal/graph` (coocorrência + export GEXF)
- [ ] Criar `internal/report/html.go` (HTML básico)
- [ ] Adicionar worker Python (`/internal/worker`) com `/lemmatize`
- [ ] Escrever testes unitários mínimos
- [ ] Docker Compose (api + worker + redis)
- [ ] Atualizar README e docs de uso

> 🔗 ver também: `.sapientia/TODO_STARTUP.md`

---

## 🧠 Orientações para agentes

Quando um agente ler este arquivo, ele deve:

1. Entender em **qual versão** o time está (atual: `v0.1`).
2. Priorizar **tarefas da sprint atual**.
3. Sugerir código, testes e docs que **avancem o estado do roadmap**.
4. Evitar propor features de versões futuras (v0.2+), a menos que explicitamente solicitado.
5. Atualizar o status dos itens conforme entregas forem feitas (ex: 🚧 → ✅).

---

## 🪄 Convenções de status

| Emoji | Significado             |
| ----- | ----------------------- |
| ✅    | Concluído               |
| 🚧    | Em desenvolvimento      |
| ⏳    | Planejado               |
| 🌙    | Longo prazo / visão     |
| 🧠    | Pesquisa ou experimento |

---

## 💬 Como atualizar o roadmap

1. Atualize o progresso manualmente nos checkboxes.
2. Mantenha a data de atualização no topo.
3. Se o roadmap mudar significativamente, avise no `README.md` principal.

---

## 📎 Referências internas

- `.sapientia/CHATMODE_SAPIENTIA.md` — comportamento e mentalidade do agente
- `.sapientia/TODO_STARTUP.md` — sprint detalhada atual
- `.sapientia/CONTEXT_ARCHITECTURE.md` — arquitetura detalhada
- `.github/CONTEXT_PROJECT.md` — visão pública
