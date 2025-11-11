package tree

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder_BuildFromLexical(t *testing.T) {
	// Criar dados de teste
	testData := LexicalData{
		Freq: map[string]int{
			"termo1": 10,
			"termo2": 5,
			"termo3": 3,
		},
		TFIDF: map[string]float64{
			"termo1": 0.5,
			"termo2": 0.3,
			"termo3": 0.1,
		},
	}

	// Criar arquivo temporário
	tmpDir := t.TempDir()
	lexicalPath := filepath.Join(tmpDir, "lexical.json")
	data, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("erro ao criar dados de teste: %v", err)
	}
	if err := os.WriteFile(lexicalPath, data, 0644); err != nil {
		t.Fatalf("erro ao escrever arquivo de teste: %v", err)
	}

	// Testar construção
	builder := NewBuilder(3)
	tree, err := builder.BuildFromLexical(lexicalPath)
	if err != nil {
		t.Fatalf("erro ao construir árvore: %v", err)
	}

	// Validar árvore
	if tree.Root == nil {
		t.Fatal("root não deveria ser nil")
	}
	if tree.Root.Term != "termo1" {
		t.Errorf("esperava root.Term = 'termo1', got %s", tree.Root.Term)
	}
	if tree.Root.Weight != 0.5 {
		t.Errorf("esperava root.Weight = 0.5, got %f", tree.Root.Weight)
	}
	if tree.MaxDepth != 3 {
		t.Errorf("esperava MaxDepth = 3, got %d", tree.MaxDepth)
	}
}

func TestBuilder_BuildFromLexicalNoTFIDF(t *testing.T) {
	// Criar dados sem TF-IDF (apenas frequências)
	testData := LexicalData{
		Freq: map[string]int{
			"termo1": 100,
			"termo2": 50,
			"termo3": 25,
		},
	}

	tmpDir := t.TempDir()
	lexicalPath := filepath.Join(tmpDir, "lexical.json")
	data, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("erro ao criar dados de teste: %v", err)
	}
	if err := os.WriteFile(lexicalPath, data, 0644); err != nil {
		t.Fatalf("erro ao escrever arquivo de teste: %v", err)
	}

	// Construir árvore
	builder := NewBuilder(2)
	tree, err := builder.BuildFromLexical(lexicalPath)
	if err != nil {
		t.Fatalf("erro ao construir árvore: %v", err)
	}

	// Validar normalização de frequências
	if tree.Root == nil {
		t.Fatal("root não deveria ser nil")
	}
	if tree.Root.Term != "termo1" {
		t.Errorf("termo com maior frequência deveria ser root")
	}
	if tree.Root.Weight != 1.0 {
		t.Errorf("peso normalizado deveria ser 1.0, got %f", tree.Root.Weight)
	}
}

func TestTree_SaveJSON(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Term:   "teste",
			Weight: 0.5,
			Level:  0,
			Children: []*Node{
				{Term: "filho1", Weight: 0.3, Level: 1},
				{Term: "filho2", Weight: 0.2, Level: 1},
			},
		},
		MaxDepth: 2,
		Metadata: map[string]interface{}{
			"test": true,
		},
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "tree.json")

	if err := tree.SaveJSON(outputPath); err != nil {
		t.Fatalf("erro ao salvar JSON: %v", err)
	}

	// Verificar arquivo
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("arquivo JSON não foi criado")
	}

	// Ler e validar conteúdo
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("erro ao ler JSON: %v", err)
	}

	var loaded Tree
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("erro ao parsear JSON: %v", err)
	}

	if loaded.Root.Term != "teste" {
		t.Errorf("esperava root.Term = 'teste', got %s", loaded.Root.Term)
	}
	if len(loaded.Root.Children) != 2 {
		t.Errorf("esperava 2 filhos, got %d", len(loaded.Root.Children))
	}
}

func TestRenderer_RenderSVG(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Term:   "raiz",
			Weight: 1.0,
			Level:  0,
			Children: []*Node{
				{
					Term:   "filho1",
					Weight: 0.5,
					Level:  1,
					Children: []*Node{
						{Term: "neto1", Weight: 0.3, Level: 2},
					},
				},
				{Term: "filho2", Weight: 0.4, Level: 1},
			},
		},
		MaxDepth: 3,
	}

	tmpDir := t.TempDir()
	svgPath := filepath.Join(tmpDir, "tree.svg")

	renderer := NewRenderer(tree)
	if err := renderer.RenderSVG(svgPath); err != nil {
		t.Fatalf("erro ao renderizar SVG: %v", err)
	}

	// Verificar arquivo
	if _, err := os.Stat(svgPath); os.IsNotExist(err) {
		t.Fatal("arquivo SVG não foi criado")
	}

	// Verificar conteúdo básico
	data, err := os.ReadFile(svgPath)
	if err != nil {
		t.Fatalf("erro ao ler SVG: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Fatal("SVG está vazio")
	}

	// Verificar tags essenciais
	required := []string{"<svg", "<circle", "<text", "</svg>"}
	for _, tag := range required {
		if !contains(content, tag) {
			t.Errorf("SVG não contém tag esperada: %s", tag)
		}
	}
}

func TestRenderer_EmptyTree(t *testing.T) {
	tree := &Tree{
		Root:     nil,
		MaxDepth: 3,
	}

	tmpDir := t.TempDir()
	svgPath := filepath.Join(tmpDir, "empty.svg")

	renderer := NewRenderer(tree)
	err := renderer.RenderSVG(svgPath)
	if err == nil {
		t.Fatal("esperava erro ao renderizar árvore vazia")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) &&
		findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
