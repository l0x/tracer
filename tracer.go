package main

import (
	"penrodyn.com/tracer/internal/ppm"
	"penrodyn.com/tracer/internal/vec"
)

func main() {
	cam := Cam{
		pos:        vec.Vec3{X: 0, Y: 2, Z: -4},
		pointingAt: vec.Vec3{X: 0, Y: 0, Z: 0},
		up:         vec.Vec3{X: 0, Y: 1, Z: 0},
		xRes:       800,
		yRes:       600,
		canvasDist: 1.0,
	}

	purple := Material{
		specularCoeff: 30.0,
		reflectivity:  0.0,
		metallic:      0.2,
		colour:        vec.Vec3{0.05, 0.0, 0.03},
	}
	green := Material{
		specularCoeff: 400.0,
		reflectivity:  0.0,
		colour:        vec.Vec3{0.8, 0.8, 0.8},
		metallic:      0.4,
	}

	floor := Plane{-1.0, &purple}
	sphere := Sphere{vec.Vec3{X: 0.0, Y: 0.2, Z: 0.0}, 0.65, &green}

	keyLight := Light{
		position:   vec.Vec3{X: 6, Y: 8, Z: -6},
		brightness: 150,
	}
	fillLight := Light{
		position:   vec.Vec3{X: -5, Y: 0, Z: -5},
		brightness: 8,
	}

	scene := Scene{
		objects:    []Visible{&floor, &sphere},
		lights:     []Light{keyLight, fillLight},
		background: vec.Vec3{0.0, 0.0, 0.0},
	}

	i := Render(&cam, &scene)
	ppm.WritePPM(i, "out.ppm")
}
