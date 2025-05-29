package main

import (
	"math"

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
	// Hybrid analytical and stepping ray tracing.
	// We need to do both, and return the nearest intersection
	// We could optimise further by tracking the intersection distance more smartly than calculating it just in time

	// First check the analytical prims
	distance := math.MaxFloat64              // nearest analytical hit distance
	var analyticalIntersection *Intersection // nearest intersection info

	for i := 0; i < len(scene.analytic); i++ {
		if hit, intrsct := scene.analytic[i].Intersect(ray); hit {
			if d := ray.pos.Sub(intrsct.pos).Magnitude(); d < distance {
				distance = d
				analyticalIntersection = intrsct
			}
		}
	}

	// Now step the ray for everything else
	rayDist := 0.0
	stepDist := ray.step.Magnitude()

	for i := 0; i < MAX_STEPS; i++ {
		ray.Step()
		rayDist += stepDist

		if rayDist > distance {
			// We already have a nearer analytical hit
			// set the ray position as if we stepped it
			ray.pos = analyticalIntersection.pos
			return analyticalIntersection
		}

		// check for collisions and return intersection
		for o := 0; o < len(scene.objects); o++ {
			if hit, intersection := scene.objects[o].Intersects(ray); hit {
				return intersection
			}
		}
	}

	return nil
}
