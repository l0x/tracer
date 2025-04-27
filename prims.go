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
	extent, origin vec.Vec3
	material       *Material
}

func NewCuboid(origin, dims vec.Vec3, material *Material) *Cuboid {
	return &Cuboid{
		origin:   origin,
		extent:   origin.Add(dims),
		material: material,
	}
}

func (c *Cuboid) Intersects(ray *Ray) (bool, *Intersection) {
	p := ray.pos

	if p.X >= c.origin.X && p.X <= c.extent.X &&
		p.Y >= c.origin.Y && p.Y <= c.extent.Y &&
		p.Z >= c.origin.Z && p.Z <= c.extent.Z {

		return true, &Intersection{
			normal:   vec.Vec3{0, 0, 0}, // TODO: WHATS THE NORMAL KENNETH?
			pos:      p,
			material: c.material,
		}
	}

	return false, nil
}
