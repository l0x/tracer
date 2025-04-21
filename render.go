package main

import (
	"math"

	"penrodyn.com/tracer/internal/img"
	"penrodyn.com/tracer/internal/vec"
)

const MAX_STEPS = 1500

func Render(cam *Cam, scene *Scene) *img.Img {
	image := img.NewImg(cam.xRes, cam.yRes)

	// Initialise the canvas location in world-space
	camPointDir := cam.pointingAt.Sub(cam.pos).Norm()

	left := cam.up.Cross(camPointDir).Norm()
	up := camPointDir.Cross(left).Norm()

	canvasCenter := cam.pos.Add(camPointDir.Scale(cam.canvasDist))
	canvasHeight := 1.0
	canvasWidth := (float64(cam.xRes) / float64(cam.yRes)) * canvasHeight

	yStep := up.Scale((canvasHeight * -1) / float64(cam.yRes))  // -1, so down
	xStep := left.Scale((canvasWidth * -1) / float64(cam.xRes)) // -1, so right

	// start in the top left
	p := canvasCenter.Add(up.Scale(canvasHeight / 2.0)).Add(left.Scale(canvasWidth / 2.0))

	// Fire a Ray through each pixel
	for y := 0; y < cam.yRes; y++ {
		for x := 0; x < cam.xRes; x++ {
			ray := NewRay(p, p.Sub(cam.pos).Norm())

			var pxVal *vec.Vec3

			intersection := ray.Trace(scene)
			if intersection != nil {
				pxVal = calcLighting(scene, ray, intersection)
			} else {
				pxVal = scene.GetBackground(ray)
			}

			// Apply gamma correction
			pxVal.PowInPlace(1.0 / 2.2)
			image.SetPixel(x, y, uint8(pxVal.X*255), uint8(pxVal.Y*255), uint8(pxVal.Z*255))

			// Step in X
			p.Translate(xStep)
		}
		// Step in Y and reset X
		p.Translate(xStep.Scale(float64(cam.xRes) * -1.0).Add(yStep))
	}
	return image
}

func calcLighting(scene *Scene, ray *Ray, intersection *Intersection) *vec.Vec3 {
	matCol := intersection.material.colour

	// "Ambient" lighting fudge
	col := matCol.Scale(0.007)

	for l := 0; l < len(scene.lights); l++ {
		// Light properties
		light := scene.lights[l]
		toLight := light.position.Sub(intersection.pos)
		lightDir := toLight.Norm()

		// Check for shadow
		shadowRay := NewRay(intersection.pos.Add(intersection.normal.Scale(0.01)), lightDir)
		if shadowRay.Trace(scene) != nil {
			continue
		}

		lightDistance := toLight.Magnitude()
		attenuation := 1.0 / (lightDistance * lightDistance)
		intensity := light.brightness * attenuation

		// Diffuse (lambertian) contribution
		diffuse := math.Max(intersection.normal.Dot(lightDir), 0) * intensity

		// Specular
		halfAngle := lightDir.Add(ray.direction.Scale(-1.0)).Norm()
		specular := math.Pow(
			math.Max(halfAngle.Dot(intersection.normal), 0),
			intersection.material.specularCoeff,
		) * intensity

		// combine the components
		kS := intersection.material.metallic // Specular factor
		kD := 1.0 - kS                       // Diffuse factor
		illumination := specular*kS + diffuse*kD

		col = col.Add(matCol.Scale(illumination))
	}

	col.Clamp()
	return &col
}
