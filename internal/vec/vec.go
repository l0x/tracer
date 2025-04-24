package vec

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func (a Vec3) Add(b Vec3) Vec3 {
	return Vec3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vec3) Sub(b Vec3) Vec3 {
	return Vec3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (v Vec3) Scale(scalar float64) Vec3 {
	return Vec3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

// Pop pop!
func (v Vec3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Norm() Vec3 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vec3{0, 0, 0}
	}
	return v.Scale(1 / mag)
}

func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Vec3) Dot(b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vec3) Lerp(b Vec3, mix float64) Vec3 {
	return Vec3{
		X: a.X*(1.0-mix) + (b.X * mix),
		Y: a.Y*(1.0-mix) + (b.Y * mix),
		Z: a.Z*(1.0-mix) + (b.Z * mix),
	}
}

func (v *Vec3) Translate(by Vec3) {
	v.X += by.X
	v.Y += by.Y
	v.Z += by.Z
}

func (v *Vec3) Clamp() {
	if v.X > 1.0 {
		v.X = 1.0
	}
	if v.Y > 1.0 {
		v.Y = 1.0
	}
	if v.Z > 1.0 {
		v.Z = 1.0
	}
}

func (v *Vec3) PowInPlace(exp float64) {
	v.X = math.Pow(v.X, exp)
	v.Y = math.Pow(v.Y, exp)
	v.Z = math.Pow(v.Z, exp)
}

func (v Vec3) Reflect(normal Vec3) Vec3 {
	return v.Sub(normal.Scale(2 * v.Dot(normal)))
}
