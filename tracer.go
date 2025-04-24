package main

import (
	"fmt"

	"penrodyn.com/tracer/internal/ppm"
	"penrodyn.com/tracer/internal/vec"
)

func main() {
	cam := Cam{
		pos:        vec.Vec3{X: 0, Y: 3, Z: -4},
		pointingAt: vec.Vec3{X: 0, Y: 0, Z: 0},
		up:         vec.Vec3{X: 0, Y: 1, Z: 0},
		xRes:       800,
		yRes:       600,
		canvasDist: 1.0,
	}

	purple := Material{
		glossy:       0.2,
		reflectivity: 0.0,
		metallic:     0.0,
		colour:       vec.Vec3{0.5, 0.0, 0.3},
	}

	whitePlastic := Material{
		glossy:       0.6,
		reflectivity: 0.1,
		metallic:     0.3,
		colour:       vec.Vec3{0.8, 0.8, 0.8},
	}

	chrome := Material{
		glossy:       0.95,
		reflectivity: 0.92,
		metallic:     0.95,
		colour:       vec.Vec3{0.8, 0.8, 0.8},
	}

	gold := Material{
		glossy:       0.2,
		reflectivity: 0.0,
		metallic:     0.96,
		colour:       vec.Vec3{0.9, 0.6, 0.05},
	}

	// seed := time.Now().UnixNano()
	// var seed int64 = 1745249485663317000

	floor := Plane{-1.0, &purple}
	sceneObjects := []Visible{&floor}

	sceneObjects = append(
		sceneObjects,
		&Sphere{
			origin:   vec.Vec3{X: 0, Y: -0.6, Z: -0.7},
			radius:   0.4,
			material: &whitePlastic,
		},
		&Sphere{
			origin:   vec.Vec3{X: -1.0, Y: -0.4, Z: -0.1},
			radius:   0.6,
			material: &chrome,
		},
		&Sphere{
			origin:   vec.Vec3{X: 0.8, Y: -0.2, Z: 0.6},
			radius:   0.8,
			material: &gold,
		},
	)

	keyLight := PointLight{
		position:   vec.Vec3{X: -3, Y: 3, Z: -5},
		brightness: 40,
	}
	fillLight := PointLight{
		position:   vec.Vec3{X: 2, Y: 1, Z: -3},
		brightness: 3,
	}

	scene := Scene{
		objects:    sceneObjects,
		lights:     []PointLight{keyLight, fillLight},
		background: vec.Vec3{0.0, 0.0, 0.0},
	}

	i := Render(&cam, &scene)
	ppm.WritePPM(i, fmt.Sprintf("out-%d.ppm", 1234))
}

// func getRandomSpheres(seed int64) []Visible {

// 	fmt.Println("Sphere seed", seed)
// 	source := rand.NewSource(seed)
// 	r := rand.New(source)

// 	spheres := []Visible{}

// 	for i := 0; i < 10; i++ {
// 		spheres = append(
// 			spheres,
// 			&Sphere{
// 				origin: vec.Vec3{
// 					X: (r.Float64() - 0.5) * 6,
// 					Y: 0.5 + (r.Float64()-0.5)*2,
// 					Z: ((r.Float64() - 0.5) * 5) + 2,
// 				},
// 				radius: 0.2 + (r.Float64() * 0.4),
// 				material: &Material{
// 					specularCoeff: 1000.0,
// 					reflectivity:  0.7,
// 					colour:        vec.Vec3{r.Float64() * 0.8, r.Float64() * 0.8, r.Float64() * 0.8},
// 					metallic:      0.9,
// 				},
// 			},
// 		)
// 	}

// 	return spheres
// }
