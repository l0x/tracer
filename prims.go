package main

import (
	"math"

	"penrodyn.com/tracer/internal/vec"
)

type Sphere struct {
	origin   vec.Vec3
	radius   float64
	material *Material
}

func (s *Sphere) Intersects(ray *Ray) (bool, *Intersection) {
	if fromCenter := ray.pos.Sub(s.origin); fromCenter.Magnitude() <= s.radius {
		return true, &Intersection{fromCenter.Norm(), ray.pos, s.material}
	}
	return false, nil
}

type FloorPlane struct {
	height   float64
	material *Material
}

func (p *FloorPlane) Intersects(ray *Ray) (bool, *Intersection) {
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

func (c *Cuboid) Intersect(ray *Ray) (bool, *Intersection) {
	// Analytic, no stepping required, assume ray.pos is origin

	var tx1, tx2, ty1, ty2, tz1, tz2 float64

	// X Slab
	if ray.direction.X > 0 {
		tx1 = (c.origin.X - ray.pos.X) / ray.direction.X
		tx2 = (c.extent.X - ray.pos.X) / ray.direction.X
	} else {
		tx2 = (c.origin.X - ray.pos.X) / ray.direction.X
		tx1 = (c.extent.X - ray.pos.X) / ray.direction.X
	}

	// Y Slab
	if ray.direction.Y > 0 {
		ty1 = (c.origin.Y - ray.pos.Y) / ray.direction.Y
		ty2 = (c.extent.Y - ray.pos.Y) / ray.direction.Y
	} else {
		ty2 = (c.origin.Y - ray.pos.Y) / ray.direction.Y
		ty1 = (c.extent.Y - ray.pos.Y) / ray.direction.Y
	}

	// Z Slab
	if ray.direction.Z > 0 {
		tz1 = (c.origin.Z - ray.pos.Z) / ray.direction.Z
		tz2 = (c.extent.Z - ray.pos.Z) / ray.direction.Z
	} else {
		tz2 = (c.origin.Z - ray.pos.Z) / ray.direction.Z
		tz1 = (c.extent.Z - ray.pos.Z) / ray.direction.Z
	}

	tmin := min(tx1, ty1, tz1)
	tmax := max(tx2, ty2, tz2)

	if tmax < 0 {
		// Box behind ray
		return false, nil
	}

	if tmin > tmax {
		// Missed the box
		return false, nil
	}

	// The normal is the face normal of the dominant axis
	var norm vec.Vec3
	if tmin == tx1 {
		if ray.direction.X > 0 {
			norm = vec.Vec3{X: -1, Y: 0, Z: 0}
		} else {
			norm = vec.Vec3{X: 1, Y: 0, Z: 0}
		}
	} else if tmin == ty1 {
		if ray.direction.Y > 0 {
			norm = vec.Vec3{X: 0, Y: -1, Z: 0}
		} else {
			norm = vec.Vec3{X: 0, Y: 1, Z: 0}
		}
	} else if tmin == tz1 {
		if ray.direction.Z > 0 {
			norm = vec.Vec3{X: 0, Y: 0, Z: -1}
		} else {
			norm = vec.Vec3{X: 0, Y: 0, Z: 1}
		}
	}

	return true, &Intersection{
		normal:   norm,
		pos:      ray.pos,
		material: c.material,
	}
}

func (s *Sphere) Intersect(ray *Ray) (bool, *Intersection) {
	oc := ray.pos.Sub(s.origin)
	a := ray.direction.Dot(ray.direction)
	b := 2 * oc.Dot(ray.direction)
	c := oc.Dot(oc) - s.radius*s.radius

	d := b*b - 4*a*c

	if d < 0 {
		return false, nil // no intersection
	}

	sqrtD := math.Sqrt(d)
	t1 := (-b - sqrtD) / (2 * a)
	t2 := (-b + sqrtD) / (2 * a)

	// Sort t1 and t2 so t1 is smaller
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	// Choose the smallest positive t
	var t float64
	if t1 > 1e-6 {
		t = t1
	} else if t2 > 1e-6 {
		t = t2
	} else {
		return false, nil // both intersections are behind the ray origin
	}

	hitPoint := ray.pos.Add(ray.direction.Scale(t))
	normal := hitPoint.Sub(s.origin).Norm()

	return true, &Intersection{
		pos:      hitPoint,
		normal:   normal,
		material: s.material,
	}
}

// Plane intersection
func (p *FloorPlane) Intersect(ray *Ray) (bool, *Intersection) {
	if ray.direction.Y >= 0 {
		return false, nil // The ray is pointing up, floors can only be seen from above
	}

	rayDist := (ray.pos.Y - p.height) / ray.direction.Y * -1
	hitPoint := ray.pos.Add(ray.direction.Scale(rayDist))

	return true, &Intersection{
		pos:      hitPoint,
		normal:   vec.Vec3{X: 0, Y: 1, Z: 0},
		material: p.material,
	}
}
