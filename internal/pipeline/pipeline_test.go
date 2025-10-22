package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPipeline(t *testing.T) {
	// Criar arquivo temporário de pipeline
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	
	yamlContent := `name: test-pipeline
description: Pipeline de teste
version: "1.0"
steps:
  - name: ingest
    type: ingest
    input: data/test.txt
    output: runs/test/corpus.json
config:
  output_dir: runs/test
`
	
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Erro ao criar arquivo de teste: %v", err)
	}
	
	// Testar LoadPipeline
	p, err := LoadPipeline(yamlPath)
	if err != nil {
		t.Fatalf("LoadPipeline falhou: %v", err)
	}
	
	// Verificar campos
	if p.Name != "test-pipeline" {
		t.Errorf("Nome esperado 'test-pipeline', obteve '%s'", p.Name)
	}
	
	if p.Description != "Pipeline de teste" {
		t.Errorf("Descrição esperada 'Pipeline de teste', obteve '%s'", p.Description)
	}
	
	if len(p.Steps) != 1 {
		t.Errorf("Esperado 1 step, obteve %d", len(p.Steps))
	}
	
	if len(p.Steps) > 0 {
		step := p.Steps[0]
		if step.Name != "ingest" {
			t.Errorf("Nome do step esperado 'ingest', obteve '%s'", step.Name)
		}
		if step.Type != "ingest" {
			t.Errorf("Tipo do step esperado 'ingest', obteve '%s'", step.Type)
		}
	}
}

func TestLoadPipeline_FileNotFound(t *testing.T) {
	_, err := LoadPipeline("nao-existe.yaml")
	if err == nil {
		t.Error("Esperado erro para arquivo inexistente")
	}
}

func TestValidate_Success(t *testing.T) {
	p := &Pipeline{
		Name:    "valid-pipeline",
		Version: "1.0",
		Steps: []Step{
			{
				Name:   "step1",
				Type:   "ingest",
				Input:  "data/input.txt",
				Output: "runs/output.json",
			},
			{
				Name:   "step2",
				Type:   "lexical",
				Input:  "runs/output.json",
				Output: "runs/lexical.json",
			},
		},
	}
	
	err := p.Validate()
	if err != nil {
		t.Errorf("Validate falhou inesperadamente: %v", err)
	}
}

func TestValidate_EmptyName(t *testing.T) {
	p := &Pipeline{
		Name:    "",
		Version: "1.0",
		Steps: []Step{
			{Name: "step1", Type: "ingest"},
		},
	}
	
	err := p.Validate()
	if err == nil {
		t.Error("Esperado erro para nome vazio")
	}
}

func TestValidate_NoSteps(t *testing.T) {
	p := &Pipeline{
		Name:    "no-steps",
		Version: "1.0",
		Steps:   []Step{},
	}
	
	err := p.Validate()
	if err == nil {
		t.Error("Esperado erro para pipeline sem steps")
	}
}

func TestValidate_EmptyStepName(t *testing.T) {
	p := &Pipeline{
		Name:    "invalid-step",
		Version: "1.0",
		Steps: []Step{
			{Name: "", Type: "ingest"},
		},
	}
	
	err := p.Validate()
	if err == nil {
		t.Error("Esperado erro para step sem nome")
	}
}

func TestValidate_EmptyStepType(t *testing.T) {
	p := &Pipeline{
		Name:    "invalid-type",
		Version: "1.0",
		Steps: []Step{
			{Name: "step1", Type: ""},
		},
	}
	
	err := p.Validate()
	if err == nil {
		t.Error("Esperado erro para step sem type")
	}
}

func TestValidate_InvalidStepType(t *testing.T) {
	p := &Pipeline{
		Name:    "invalid-type",
		Version: "1.0",
		Steps: []Step{
			{Name: "step1", Type: "invalid-type"},
		},
	}
	
	err := p.Validate()
	if err == nil {
		t.Error("Esperado erro para tipo de step inválido")
	}
}

func TestValidate_AllStepTypes(t *testing.T) {
	validTypes := []string{"ingest", "lexical", "graph", "tree", "report"}
	
	for _, stepType := range validTypes {
		p := &Pipeline{
			Name:    "test-" + stepType,
			Version: "1.0",
			Steps: []Step{
				{Name: "step1", Type: stepType},
			},
		}
		
		err := p.Validate()
		if err != nil {
			t.Errorf("Validate falhou para tipo válido '%s': %v", stepType, err)
		}
	}
}
