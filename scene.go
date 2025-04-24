package main

import "penrodyn.com/tracer/internal/vec"

type Intersection struct {
	normal, pos vec.Vec3
	material    *Material
}

type Material struct {
	glossy       float64
	reflectivity float64
	colour       vec.Vec3
	metallic     float64
}

type Visible interface {
	Intersects(ray *Ray) (bool, *Intersection)
}

type PointLight struct {
	position   vec.Vec3
	brightness float64
}

type Scene struct {
	objects    []Visible
	lights     []PointLight
	background vec.Vec3
}

func (scene *Scene) GetBackground(ray *Ray) *vec.Vec3 {
	// Just return a flat background colour for now
	// but we could use the ray to make a gradient...
	return &scene.background
}
