// Package pipeline implementa o executor de pipelines
package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Executor executa um pipeline completo
type Executor struct {
	pipeline *Pipeline
	workDir  string
}

// NewExecutor cria um novo executor para o pipeline
func NewExecutor(p *Pipeline) *Executor {
	return &Executor{
		pipeline: p,
		workDir:  ".",
	}
}

// Execute executa o pipeline completo
func (e *Executor) Execute() error {
	fmt.Printf("📋 Pipeline: %s\n", e.pipeline.Name)
	if e.pipeline.Description != "" {
		fmt.Printf("   %s\n", e.pipeline.Description)
	}
	fmt.Println()

	// Validar antes de executar
	if err := e.pipeline.Validate(); err != nil {
		return fmt.Errorf("validação falhou: %w", err)
	}

	// Criar diretório de output se especificado
	if outputDir, ok := e.pipeline.Config["output_dir"]; ok {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar output_dir: %w", err)
		}
		fmt.Printf("📁 Output: %s\n\n", outputDir)
	}

	// Executar steps sequencialmente
	for i, step := range e.pipeline.Steps {
		start := time.Now()
		fmt.Printf("▶️  Step %d/%d: %s (%s)\n", i+1, len(e.pipeline.Steps), step.Name, step.Type)

		if err := e.executeStep(&step); err != nil {
			return fmt.Errorf("step %s falhou: %w", step.Name, err)
		}

		elapsed := time.Since(start)
		fmt.Printf("   ✅ Concluído em %.2fs\n\n", elapsed.Seconds())
	}

	fmt.Println("🎉 Pipeline executado com sucesso!")
	return nil
}

// executeStep executa um step individual
func (e *Executor) executeStep(step *Step) error {
	switch step.Type {
	case "ingest":
		return e.executeIngest(step)
	case "lexical":
		return e.executeLexical(step)
	case "graph":
		return e.executeGraph(step)
	case "tree":
		return e.executeTree(step)
	case "report":
		return e.executeReport(step)
	default:
		return fmt.Errorf("tipo de step desconhecido: %s", step.Type)
	}
}

// executeIngest executa step de ingestão
func (e *Executor) executeIngest(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Ingest ainda não implementado (v0.1)")
	
	// Criar arquivo de output vazio como placeholder
	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		if err := os.WriteFile(step.Output, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeLexical executa análise lexical
func (e *Executor) executeLexical(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Lexical ainda não implementado (v0.1)")
	
	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		if err := os.WriteFile(step.Output, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeGraph executa geração de grafo
func (e *Executor) executeGraph(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Graph ainda não implementado (v0.1)")
	
	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		// Criar GEXF mínimo como placeholder
		gexf := `<?xml version="1.0" encoding="UTF-8"?>
<gexf xmlns="http://www.gexf.net/1.2draft" version="1.2">
  <graph mode="static" defaultedgetype="undirected">
    <nodes></nodes>
    <edges></edges>
  </graph>
</gexf>`
		if err := os.WriteFile(step.Output, []byte(gexf), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeTree executa geração de árvore semântica
func (e *Executor) executeTree(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Tree ainda não implementado (v0.1)")
	
	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		if err := os.WriteFile(step.Output, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeReport executa geração de relatório
func (e *Executor) executeReport(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Report ainda não implementado (v0.1)")
	
	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		// Criar HTML mínimo como placeholder
		html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>SapientIA Report</title>
</head>
<body>
    <h1>Relatório SapientIA</h1>
    <p>Relatório em desenvolvimento (v0.1)</p>
</body>
</html>`
		if err := os.WriteFile(step.Output, []byte(html), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// ensureOutputDir garante que o diretório de output existe
func (e *Executor) ensureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório: %w", err)
		}
	}
	return nil
}
