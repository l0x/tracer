package main

import (
	"penrodyn.com/tracer/internal/vec"
)

type Cam struct {
	pos, pointingAt, up vec.Vec3
	xRes, yRes          int
	canvasDist          float64
}
