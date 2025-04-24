package main

import (
	"math"
	"math/rand"
	"sync"

	"penrodyn.com/tracer/internal/img"
	"penrodyn.com/tracer/internal/vec"
)

const MAX_STEPS = 2000
const N_REFLECTIONS = 5
const GLOSS_RAMP = 3
const MAX_SHINY = 2000
const GLOSSY_REFLECTIONS = false

func Render(cam *Cam, scene *Scene) *img.Img {
	image := img.NewImg(cam.xRes, cam.yRes)

	// Initialise the canvas location in world-space
	camPointDir := cam.pointingAt.Sub(cam.pos).Norm()

	left := camPointDir.Cross(cam.up).Norm()
	//left := cam.up.Cross(camPointDir).Norm()
	up := left.Cross(camPointDir).Norm()
	//up := camPointDir.Cross(left).Norm()

	canvasCenter := cam.pos.Add(camPointDir.Scale(cam.canvasDist))
	canvasHeight := 1.0
	canvasWidth := (float64(cam.xRes) / float64(cam.yRes)) * canvasHeight

	yStep := up.Scale((canvasHeight * -1) / float64(cam.yRes))  // -1, so down
	xStep := left.Scale((canvasWidth * -1) / float64(cam.xRes)) // -1, so right

	// start in the top left
	tl := canvasCenter.Add(up.Scale(canvasHeight / 2.0)).Add(left.Scale(canvasWidth / 2.0))

	var wg sync.WaitGroup
	wg.Add(cam.yRes)

	// Fire a Ray through each pixel
	for y := 0; y < cam.yRes; y++ {
		go func(y int) {
			defer wg.Done()
			for x := 0; x < cam.xRes; x++ {
				p := tl.Add(yStep.Scale(float64(y)).Add(xStep.Scale(float64(x))))
				ray := NewRay(p, p.Sub(cam.pos).Norm())

				var pxVal *vec.Vec3

				intersection := ray.Trace(scene)
				if intersection != nil {
					pxVal = calcLighting(scene, ray, intersection, N_REFLECTIONS)
				} else {
					pxVal = scene.GetBackground(ray)
				}

				// Apply gamma correction
				pxVal.PowInPlace(1.0 / 2.2)
				image.SetPixel(x, y, uint8(pxVal.X*255), uint8(pxVal.Y*255), uint8(pxVal.Z*255))
			}
		}(y)
	}

	wg.Wait()
	return image
}

func calcLighting(scene *Scene, ray *Ray, intersection *Intersection, reflections int) *vec.Vec3 {
	material := intersection.material

	// "Ambient" lighting fudge
	//col := material.colour.Scale(0.007)
	col := vec.Vec3{0, 0, 0}

	pointJustOffSurface := intersection.pos.Add(intersection.normal.Scale(0.01))

	for l := 0; l < len(scene.lights); l++ {
		// Light properties
		light := scene.lights[l]
		toLight := light.position.Sub(intersection.pos)
		lightDir := toLight.Norm()

		// Check for shadow
		shadowRay := NewRay(pointJustOffSurface, lightDir)
		if shadowRay.Trace(scene) != nil {
			continue
		}

		lightDistance := toLight.Magnitude()
		attenuation := 1.0 / (lightDistance * lightDistance)
		intensity := light.brightness * attenuation

		// Diffuse (lambertian) contribution
		diffuse := math.Max(intersection.normal.Dot(lightDir), 0) * intensity

		// Specular
		phongCoefficient := math.Pow(material.glossy, GLOSS_RAMP) * MAX_SHINY
		halfAngle := lightDir.Add(ray.direction.Scale(-1.0)).Norm()
		specular := math.Pow(
			math.Max(halfAngle.Dot(intersection.normal), 0),
			phongCoefficient,
		) * intensity

		// combine the components
		kS := material.metallic // Specular factor
		kD := 1.0 - kS          // Diffuse factorr
		illumination := specular*kS + diffuse*kD

		col = col.Add(material.colour.Scale(illumination))
	}

	// Reflection
	if reflections > 0 && material.reflectivity > 0 {
		reflected := ray.direction.Reflect(intersection.normal)

		// Perturb it by using material glossyness parameter
		if GLOSSY_REFLECTIONS {
			perturbed := perturb(reflected)
			if perturbed.Dot(intersection.normal) < 0 {
				perturbed = perturbed.Scale(-1)
			}
			reflected = reflected.Norm().Lerp(perturbed.Norm(), 1.0-material.glossy).Norm()
		}

		reflectRay := NewRay(pointJustOffSurface, reflected)
		intersection := reflectRay.Trace(scene)

		var reflection *vec.Vec3
		if intersection == nil {
			reflection = scene.GetBackground(reflectRay)
		} else {
			reflection = calcLighting(scene, reflectRay, intersection, reflections-1)
		}

		// Mix it in
		col = col.Scale(1.0 - material.reflectivity)
		col = col.Add(reflection.Scale(material.reflectivity))
	}

	col.Clamp()
	return &col
}

func CosineSampleHemisphere() vec.Vec3 {
	u1 := rand.Float64()
	u2 := rand.Float64()

	r := math.Sqrt(u1)
	theta := 2 * math.Pi * u2

	x := r * math.Cos(theta)
	y := r * math.Sin(theta)
	z := math.Sqrt(1 - u1)

	return vec.Vec3{x, y, z}
}

func perturb(dir vec.Vec3) vec.Vec3 {
	local := CosineSampleHemisphere()

	// make orthonormal basis from direction
	var u vec.Vec3
	if math.Abs(dir.X) > math.Abs(dir.Z) {
		u = vec.Vec3{-dir.Y, dir.X, 0}
	} else {
		u = vec.Vec3{0, -dir.Z, dir.Y}
	}

	u = u.Norm()
	v := dir.Cross(u)
	w := dir.Norm()

	// Rotate the sample relative to the direction
	return u.Scale(local.X).Add(v.Scale(local.Y)).Add(w.Scale(local.Z))
}
