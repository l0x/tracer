package main

import "penrodyn.com/tracer/internal/vec"

type Intersection struct {
	normal, pos vec.Vec3
	material    *Material
}

type Material struct {
	specularCoeff float64
	reflectivity  float64
	colour        vec.Vec3
	metallic      float64
}

type Visible interface {
	Intersects(ray *Ray) (bool, *Intersection)
}

type Light struct {
	position   vec.Vec3
	brightness float64
}

type Scene struct {
	objects    []Visible
	lights     []Light
	background vec.Vec3
}

func (scene *Scene) GetBackground(ray *Ray) *vec.Vec3 {
	// Just return a flat background colour for now
	// but we could use the ray to make a gradient...
	return &scene.background
}

type Sphere struct {
	origin   vec.Vec3
	radius   float64
	material *Material
}

func (s *Sphere) Intersects(ray *Ray) (bool, *Intersection) {
	fromCenter := ray.pos.Sub(s.origin)
	dist := fromCenter.Magnitude()

	// if fromCenter := s.origin.Sub(ray.pos); fromCenter.Magnitude() <= s.radius {
	if dist <= s.radius {
		return true, &Intersection{fromCenter.Norm(), ray.pos, s.material}
	}
	return false, nil
}

type Plane struct {
	height   float64
	material *Material
}

func (p *Plane) Intersects(ray *Ray) (bool, *Intersection) {
	if ray.pos.Y <= p.height {
		return true, &Intersection{vec.Vec3{X: 0, Y: 1, Z: 0}, ray.pos, p.material}
	}
	return false, nil
}
