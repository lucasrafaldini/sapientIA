// Package tree implementa geração de árvore semântica
package tree

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

// Renderer renderiza árvores em diferentes formatos
type Renderer struct {
	tree         *Tree
	nodeSpacing  float64
	levelSpacing float64
	fontSize     int
	width        float64
	height       float64
}

// NewRenderer cria um novo renderer
func NewRenderer(tree *Tree) *Renderer {
	return &Renderer{
		tree:         tree,
		nodeSpacing:  120,
		levelSpacing: 80,
		fontSize:     14,
		width:        1200,
		height:       800,
	}
}

// RenderSVG gera visualização SVG da árvore
func (r *Renderer) RenderSVG(path string) error {
	if r.tree.Root == nil {
		return fmt.Errorf("árvore vazia")
	}

	var svg bytes.Buffer
	
	// Calcular dimensões necessárias
	r.calculateDimensions(r.tree.Root)
	
	// Header SVG
	svg.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
  <defs>
    <style>
      .node { fill: #4a90e2; stroke: #2e5c8a; stroke-width: 2; }
      .node-text { fill: white; font-family: Arial, sans-serif; font-size: %dpx; text-anchor: middle; }
      .edge { stroke: #888; stroke-width: 2; fill: none; }
      .weight-text { fill: #666; font-family: Arial, sans-serif; font-size: 11px; text-anchor: middle; }
    </style>
  </defs>
`, r.width, r.height, r.width, r.height, r.fontSize))

	// Renderizar árvore recursivamente
	x := r.width / 2
	y := 40.0
	r.renderNode(&svg, r.tree.Root, x, y, r.width/2)

	svg.WriteString("</svg>")

	// Salvar arquivo
	if err := os.WriteFile(path, svg.Bytes(), 0644); err != nil {
		return fmt.Errorf("erro ao salvar SVG: %w", err)
	}

	return nil
}

// renderNode renderiza um nó e seus filhos recursivamente
func (r *Renderer) renderNode(svg *bytes.Buffer, node *Node, x, y, spread float64) {
	if node == nil {
		return
	}

	radius := 30.0 + node.Weight*20
	
	// Desenhar arestas para filhos primeiro (para ficarem atrás)
	if len(node.Children) > 0 {
		childY := y + r.levelSpacing
		numChildren := len(node.Children)
		childSpread := spread / 2
		
		for i, child := range node.Children {
			offset := float64(i) - float64(numChildren-1)/2
			childX := x + offset*r.nodeSpacing
			
			// Limitar posição horizontal
			childX = math.Max(50, math.Min(r.width-50, childX))
			
			// Desenhar aresta
			svg.WriteString(fmt.Sprintf(`  <line class="edge" x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"/>
`, x, y+radius, childX, childY-radius))
			
			// Renderizar filho recursivamente
			r.renderNode(svg, child, childX, childY, childSpread)
		}
	}
	
	// Desenhar nó atual
	svg.WriteString(fmt.Sprintf(`  <circle class="node" cx="%.2f" cy="%.2f" r="%.2f"/>
`, x, y, radius))
	
	// Desenhar texto do termo (quebrar se muito longo)
	term := node.Term
	if len(term) > 15 {
		term = term[:12] + "..."
	}
	svg.WriteString(fmt.Sprintf(`  <text class="node-text" x="%.2f" y="%.2f">%s</text>
`, x, y+5, escapeXML(term)))
	
	// Desenhar peso abaixo do nó
	svg.WriteString(fmt.Sprintf(`  <text class="weight-text" x="%.2f" y="%.2f">%.3f</text>
`, x, y+radius+15, node.Weight))
}

// calculateDimensions calcula dimensões necessárias baseado na árvore
func (r *Renderer) calculateDimensions(node *Node) {
	// Contar níveis e largura máxima
	maxLevel := 0
	var countNodes func(*Node, int) int
	countNodes = func(n *Node, level int) int {
		if n == nil {
			return 0
		}
		if level > maxLevel {
			maxLevel = level
		}
		count := 1
		for _, child := range n.Children {
			count += countNodes(child, level+1)
		}
		return count
	}
	
	totalNodes := countNodes(node, 0)
	
	// Ajustar dimensões
	r.height = float64(maxLevel+1)*r.levelSpacing + 100
	r.width = math.Max(800, float64(totalNodes)*r.nodeSpacing/2)
}

// RenderPNG gera visualização PNG da árvore (requer ImageMagick/rsvg-convert)
func (r *Renderer) RenderPNG(svgPath, pngPath string) error {
	// Primeiro gerar SVG temporário se não existe
	tempSVG := svgPath
	if svgPath == "" {
		tempSVG = pngPath + ".tmp.svg"
		if err := r.RenderSVG(tempSVG); err != nil {
			return err
		}
		defer os.Remove(tempSVG)
	}

	// Tentar rsvg-convert primeiro (mais comum em Linux/Mac)
	cmd := exec.Command("rsvg-convert", tempSVG, "-o", pngPath)
	if err := cmd.Run(); err != nil {
		// Se falhar, tentar ImageMagick
		cmd = exec.Command("convert", tempSVG, pngPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao converter SVG->PNG (instale rsvg-convert ou ImageMagick): %w", err)
		}
	}

	return nil
}

// escapeXML escapa caracteres especiais para XML
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
