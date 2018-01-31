package internal

import (
	"github.com/faiface/pixel"
)

type MatrixStack struct {
	stack []pixel.Matrix
}

func (m *MatrixStack) Push(mat pixel.Matrix) {
	if len(m.stack) > 0 {
		m.stack = append(m.stack, m.Head().Chained(mat))
	} else {
		m.stack = append(m.stack, mat)
	}
}

func (m *MatrixStack) Pop() pixel.Matrix {
	mat := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	return mat
}

func (m *MatrixStack) Head() pixel.Matrix {
	return m.stack[len(m.stack)-1]
}

func (m *MatrixStack) Project(v pixel.Vec) pixel.Vec {
	return m.Head().Project(v)
}
