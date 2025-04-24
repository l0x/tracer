package main

import "penrodyn.com/tracer/internal/vec"

type Sphere struct {
	origin   vec.Vec3
	radius   float64
	material *Material
}

func (s *Sphere) Intersects(ray *Ray) (bool, *Intersection) {
	if fromCenter := s.origin.Sub(ray.pos); fromCenter.Magnitude() <= s.radius {
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

type Cuboid struct {
	height, width, depth float64
	material             *Material
}

func (c *Cuboid) Intersects(ray *Ray) (bool, *Intersection) {
	return false, nil
}
