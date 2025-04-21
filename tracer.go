package main

import (
	"fmt"
	"math/rand"

	"penrodyn.com/tracer/internal/ppm"
	"penrodyn.com/tracer/internal/vec"
)

func main() {
	cam := Cam{
		pos:        vec.Vec3{X: -1, Y: 1, Z: -5},
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
	// green := Material{
	// 	specularCoeff: 400.0,
	// 	reflectivity:  0.0,
	// 	colour:        vec.Vec3{0.8, 0.8, 0.8},
	// 	metallic:      0.4,
	// }

	//seed := time.Now().UnixNano()

	var seed int64 = 1745249485663317000
	fmt.Println("Scene seed", seed)
	source := rand.NewSource(seed)
	r := rand.New(source)

	floor := Plane{-1.0, &purple}
	sceneObjects := []Visible{&floor}

	for i := 0; i < 10; i++ {
		sceneObjects = append(
			sceneObjects,
			&Sphere{
				origin: vec.Vec3{
					X: (r.Float64() - 0.5) * 6,
					Y: 0.5 + (r.Float64()-0.5)*2,
					Z: ((r.Float64() - 0.5) * 5) + 2,
				},
				radius: 0.2 + (r.Float64() * 0.4),
				material: &Material{
					specularCoeff: 10.0,
					reflectivity:  0.0,
					colour:        vec.Vec3{r.Float64() * 0.8, r.Float64() * 0.8, r.Float64() * 0.8},
					metallic:      0.9,
				},
			},
		)
	}

	//sphere := Sphere{vec.Vec3{X: 0.0, Y: 0.2, Z: 0.0}, 0.65, &green}

	keyLight := Light{
		position:   vec.Vec3{X: 6, Y: 8, Z: -6},
		brightness: 300,
	}
	fillLight := Light{
		position:   vec.Vec3{X: -10, Y: 1, Z: -10},
		brightness: 50,
	}

	scene := Scene{
		objects:    sceneObjects,
		lights:     []Light{keyLight, fillLight},
		background: vec.Vec3{0.0, 0.0, 0.0},
	}

	i := Render(&cam, &scene)
	ppm.WritePPM(i, fmt.Sprintf("out-%d.ppm", seed))
}
