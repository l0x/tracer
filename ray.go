package main

import (
	"penrodyn.com/tracer/internal/vec"
)

const STEP_DIST = 0.01

type Ray struct {
	pos, step, direction vec.Vec3
}

func NewRay(start, direction vec.Vec3) *Ray {
	return &Ray{
		pos:       start,
		direction: direction,
		step:      direction.Scale(STEP_DIST),
	}
}

func (ray *Ray) Step() {
	ray.pos.Translate(ray.step)
}

func (ray *Ray) Trace(scene *Scene) *Intersection {
	for i := 0; i < MAX_STEPS; i++ {
		ray.Step()

		// check for collisions and return intersection
		for o := 0; o < len(scene.objects); o++ {
			if hit, intersection := scene.objects[o].Intersects(ray); hit {
				return intersection
			}
		}
	}

	return nil
}
