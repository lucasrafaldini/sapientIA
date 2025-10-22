// Package pipeline implementa a leitura e validação de pipelines YAML
package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Pipeline representa um pipeline de análise
type Pipeline struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	Steps       []Step            `yaml:"steps"`
	Config      map[string]string `yaml:"config"`
}

// Step representa uma etapa do pipeline
type Step struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Input  string                 `yaml:"input"`
	Output string                 `yaml:"output"`
	Params map[string]interface{} `yaml:"params"`
}

// LoadPipeline carrega um pipeline a partir de arquivo YAML
func LoadPipeline(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler pipeline: %w", err)
	}

	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("erro ao parsear YAML: %w", err)
	}

	return &p, nil
}

// Validate valida a estrutura do pipeline
func (p *Pipeline) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("pipeline deve ter um nome")
	}
	if len(p.Steps) == 0 {
		return fmt.Errorf("pipeline deve ter ao menos um step")
	}

	// Validar cada step
	for i, step := range p.Steps {
		if step.Name == "" {
			return fmt.Errorf("step %d: nome é obrigatório", i)
		}
		if step.Type == "" {
			return fmt.Errorf("step %s: type é obrigatório", step.Name)
		}
		// Validar tipos conhecidos
		validTypes := map[string]bool{
			"ingest": true, "lexical": true, "graph": true,
			"tree": true, "report": true,
		}
		if !validTypes[step.Type] {
			return fmt.Errorf("step %s: tipo '%s' inválido (use: ingest, lexical, graph, tree, report)", step.Name, step.Type)
		}
	}

	return nil
}
