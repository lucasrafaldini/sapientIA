// Package tree implementa geração de árvore semântica
package tree

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
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
		nodeSpacing:  150, // Reduzido de 200 para 150 (menos espaço horizontal)
		levelSpacing: 250, // Aumentado de 150 para 250 (muito mais espaço vertical)
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

	radius := 40.0 + node.Weight*30 // Aumentado para nós maiores e mais visíveis

	// Desenhar arestas para filhos primeiro (para ficarem atrás)
	if len(node.Children) > 0 {
		childY := y + r.levelSpacing
		numChildren := len(node.Children)

		// Calcular largura total necessária para os filhos
		totalWidth := float64(numChildren-1) * r.nodeSpacing
		startX := x - totalWidth/2

		for i, child := range node.Children {
			childX := startX + float64(i)*r.nodeSpacing

			// Limitar posição horizontal com margem
			childX = math.Max(80, math.Min(r.width-80, childX))

			// Desenhar aresta
			fmt.Fprintf(svg, `  <line class="edge" x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"/>
`, x, y+radius, childX, childY-radius)

			// Renderizar filho recursivamente
			r.renderNode(svg, child, childX, childY, r.nodeSpacing)
		}
	}
	// Desenhar nó atual
	fmt.Fprintf(svg, `  <circle class="node" cx="%.2f" cy="%.2f" r="%.2f"/>
`, x, y, radius)

	// Desenhar texto do termo (quebrar se muito longo)
	term := node.Term
	if len(term) > 15 {
		term = term[:12] + "..."
	}
	fmt.Fprintf(svg, `  <text class="node-text" x="%.2f" y="%.2f">%s</text>
`, x, y+5, escapeXML(term))

	// Desenhar peso abaixo do nó
	fmt.Fprintf(svg, `  <text class="weight-text" x="%.2f" y="%.2f">%.3f</text>
`, x, y+radius+15, node.Weight)
}

// calculateDimensions calcula dimensões necessárias baseado na árvore
func (r *Renderer) calculateDimensions(node *Node) {
	// Contar níveis e nós por nível
	maxLevel := 0
	nodesPerLevel := make(map[int]int)

	var countNodes func(*Node, int)
	countNodes = func(n *Node, level int) {
		if n == nil {
			return
		}
		if level > maxLevel {
			maxLevel = level
		}
		nodesPerLevel[level]++
		for _, child := range n.Children {
			countNodes(child, level+1)
		}
	}

	countNodes(node, 0)

	// Calcular largura máxima necessária baseada no nível com mais nós
	maxNodesInLevel := 0
	for _, count := range nodesPerLevel {
		if count > maxNodesInLevel {
			maxNodesInLevel = count
		}
	}

	// Ajustar dimensões priorizando altura (layout vertical)
	r.height = float64(maxLevel+1)*r.levelSpacing + 200                  // Aumentada margem vertical
	r.width = math.Max(1000, float64(maxNodesInLevel)*r.nodeSpacing+150) // Reduzida largura
}

// RenderPNG gera visualização PNG da árvore usando stdlib do Go
func (r *Renderer) RenderPNG(pngPath string) error {
	if r.tree.Root == nil {
		return fmt.Errorf("árvore vazia")
	}

	// Calcular dimensões
	r.calculateDimensions(r.tree.Root)

	// Criar imagem
	img := image.NewRGBA(image.Rect(0, 0, int(r.width), int(r.height)))

	// Background branco
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Renderizar árvore
	x := r.width / 2
	y := r.height - 40.0
	r.renderNodeImage(img, r.tree.Root, x, y, r.width/2)

	// Salvar PNG
	f, err := os.Create(pngPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo PNG: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("erro ao fechar arquivo PNG: %w", cerr)
		}
	}()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("erro ao codificar PNG: %w", err)
	}

	return nil
}

// renderNodeImage renderiza um nó usando image/draw
func (r *Renderer) renderNodeImage(img *image.RGBA, node *Node, x, y, spread float64) {
	if node == nil {
		return
	}

	radius := 40.0 + node.Weight*30 // Aumentado para nós maiores e mais visíveis

	// Desenhar filhos primeiro (para ficarem atrás)
	if len(node.Children) > 0 {
		childY := y - r.levelSpacing
		numChildren := len(node.Children)

		// Calcular largura total necessária para os filhos
		totalWidth := float64(numChildren-1) * r.nodeSpacing
		startX := x - totalWidth/2

		for i, child := range node.Children {
			childX := startX + float64(i)*r.nodeSpacing

			// Limitar posição horizontal com margem
			childX = math.Max(80, math.Min(r.width-80, childX))

			// Desenhar linha
			r.drawLine(img, int(x), int(y-radius), int(childX), int(childY+radius), color.RGBA{136, 136, 136, 255})

			// Renderizar filho recursivamente
			r.renderNodeImage(img, child, childX, childY, r.nodeSpacing)
		}
	}
	// Desenhar círculo do nó
	r.drawCircle(img, int(x), int(y), int(radius), color.RGBA{74, 144, 226, 255}, color.RGBA{46, 92, 138, 255})

	// Desenhar texto do termo
	term := node.Term
	if len(term) > 15 {
		term = term[:12] + "..."
	}
	r.drawText(img, int(x), int(y), term, color.White)

	// Desenhar peso
	weightText := fmt.Sprintf("%.3f", node.Weight)
	r.drawText(img, int(x), int(y)-int(radius)-10, weightText, color.RGBA{136, 136, 136, 255})
}

// drawCircle desenha um círculo preenchido com borda
func (r *Renderer) drawCircle(img *image.RGBA, cx, cy, radius int, fill, stroke color.Color) {
	// Desenhar círculo preenchido
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				img.Set(cx+x, cy+y, fill)
			}
		}
	}

	// Desenhar borda
	thickness := 2
	for y := -radius - thickness; y <= radius+thickness; y++ {
		for x := -radius - thickness; x <= radius+thickness; x++ {
			dist := x*x + y*y
			if dist > (radius-thickness)*(radius-thickness) && dist <= (radius+thickness)*(radius+thickness) {
				img.Set(cx+x, cy+y, stroke)
			}
		}
	}
}

// drawLine desenha uma linha usando algoritmo de Bresenham
func (r *Renderer) drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// drawText desenha texto centralizado usando fonte TrueType com suporte UTF-8
func (r *Renderer) drawText(img *image.RGBA, x, y int, text string, col color.Color) {
	// Parse da fonte TrueType
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return // Silenciosamente falha se não conseguir carregar fonte
	}

	// Criar face com tamanho 13
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 13,
		DPI:  72,
	})
	defer func() {
		_ = face.Close() // Ignora erro de Close() propositalmente
	}()

	// Calcular largura do texto para centralizar
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}

	textBounds := d.MeasureString(text)
	textWidth := textBounds.Round()

	point := fixed.Point26_6{
		X: fixed.Int26_6(x*64) - fixed.Int26_6(textWidth*32),
		Y: fixed.Int26_6(y * 64),
	}

	d.Dot = point
	d.DrawString(text)
}

// abs retorna o valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
