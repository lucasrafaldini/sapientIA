// Package tree implementa geração de árvore semântica
package tree

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Node representa um nó na árvore semântica
type Node struct {
	Term     string   `json:"term"`
	Weight   float64  `json:"weight"`
	Level    int      `json:"level"`
	Children []*Node  `json:"children,omitempty"`
}

// Tree representa uma árvore semântica hierárquica
type Tree struct {
	Root     *Node                  `json:"root"`
	MaxDepth int                    `json:"max_depth"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Builder constrói árvores semânticas a partir de dados lexicais
type Builder struct {
	maxDepth   int
	minWeight  float64
	maxBranch  int
}

// NewBuilder cria um novo builder de árvores
func NewBuilder(maxDepth int) *Builder {
	return &Builder{
		maxDepth:  maxDepth,
		minWeight: 0.001, // Permitir termos com peso muito baixo
		maxBranch: 10,
	}
}

// LexicalData representa os dados de entrada do step lexical
type LexicalData struct {
	Freq     map[string]int         `json:"freq"`
	TFIDF    map[string]float64     `json:"tfidf,omitempty"`
	NGrams2  map[string]int         `json:"ngrams_2,omitempty"`
	NGrams3  map[string]int         `json:"ngrams_3,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BuildFromLexical constrói árvore a partir de dados lexicais
func (b *Builder) BuildFromLexical(lexicalPath string) (*Tree, error) {
	// Ler dados lexicais
	data, err := os.ReadFile(lexicalPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler lexical.json: %w", err)
	}

	var lexical LexicalData
	if err := json.Unmarshal(data, &lexical); err != nil {
		return nil, fmt.Errorf("erro ao parsear lexical.json: %w", err)
	}

	// Construir hierarquia
	tree := &Tree{
		MaxDepth: b.maxDepth,
		Metadata: map[string]interface{}{
			"algorithm": "tfidf-hierarchical",
			"source":    lexicalPath,
		},
	}

	// Usar TF-IDF se disponível, senão usar frequências
	weights := lexical.TFIDF
	if len(weights) == 0 {
		// Normalizar frequências
		weights = make(map[string]float64)
		maxFreq := 0
		for _, freq := range lexical.Freq {
			if freq > maxFreq {
				maxFreq = freq
			}
		}
		for term, freq := range lexical.Freq {
			weights[term] = float64(freq) / float64(maxFreq)
		}
	}

	// Criar nó raiz com termo mais importante
	tree.Root = b.buildHierarchy(weights, 0)

	return tree, nil
}

// buildHierarchy constrói hierarquia recursivamente
func (b *Builder) buildHierarchy(weights map[string]float64, level int) *Node {
	if level >= b.maxDepth || len(weights) == 0 {
		return nil
	}

	// Ordenar termos por peso
	type termWeight struct {
		term   string
		weight float64
	}
	var items []termWeight
	for term, weight := range weights {
		if weight >= b.minWeight {
			items = append(items, termWeight{term, weight})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].weight > items[j].weight
	})

	if len(items) == 0 {
		return nil
	}

	// Criar nó com termo mais relevante
	root := &Node{
		Term:   items[0].term,
		Weight: items[0].weight,
		Level:  level,
	}

	// Se não chegamos na profundidade máxima, criar filhos
	if level < b.maxDepth-1 && len(items) > 1 {
		// Particionar termos restantes em grupos
		remaining := items[1:]
		numGroups := min(b.maxBranch, len(remaining))
		
		if numGroups > 0 {
			groupSize := len(remaining) / numGroups
			if groupSize < 1 {
				groupSize = 1
			}

			for i := 0; i < numGroups && i*groupSize < len(remaining); i++ {
				start := i * groupSize
				end := start + groupSize
				if i == numGroups-1 {
					end = len(remaining)
				}
				
				// Criar subárvore para cada grupo
				subWeights := make(map[string]float64)
				for _, item := range remaining[start:end] {
					subWeights[item.term] = item.weight
				}
				
				if child := b.buildHierarchy(subWeights, level+1); child != nil {
					root.Children = append(root.Children, child)
				}
			}
		}
	}

	return root
}

// SaveJSON salva árvore em formato JSON
func (t *Tree) SaveJSON(path string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}

// min retorna o menor de dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
