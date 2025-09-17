package Engine

import (
	"math"
)

type Vector2 struct {
	X, Y float64
}

func DistanceSquare_Vector2(const_a, const_b *Vector2) float64 {
	return math.Pow(const_a.X-const_b.X, 2) + math.Pow(const_a.Y-const_b.Y, 2)
}

func Magnitude_Vector2(const_a *Vector2) float64 {
	return math.Sqrt(const_a.X*const_a.X + const_a.Y*const_a.Y)
}

func MagnitudeSquare_Vector2(const_a *Vector2) float64 {
	return (const_a.X*const_a.X + const_a.Y*const_a.Y)
}

func Normalise_Vector2(const_a *Vector2) Vector2 {
	magnitude := Magnitude_Vector2(const_a)
	return Vector2{const_a.X / magnitude, const_a.Y / magnitude}
}

func Add_Vector2(const_a, const_b *Vector2) Vector2 {
	return Vector2{const_a.X + const_b.X, const_a.Y + const_b.Y}
}

func Reflect_Vector2(const_a, const_b *Vector2) Vector2 {
	// make sure n is normalized before using this!
	dot := Dot_Vector2(const_a, const_b)
	return Vector2{
		const_a.X - 2*dot*const_b.X,
		const_a.Y - 2*dot*const_b.Y,
	}
}

func Subtract_Vector2(const_a, const_b *Vector2) Vector2 {
	return Vector2{const_a.X - const_b.X, const_a.Y - const_b.Y}
}

func Add_Float_Vector2(const_a float64, const_b *Vector2) Vector2 {
	return Vector2{const_b.X + const_a, const_b.Y + const_a}
}

func Multiply_Float_Vector2(const_a float64, const_b *Vector2) Vector2 {
	return Vector2{const_b.X * const_a, const_b.Y * const_a}
}

func Dot_Vector2(const_a *Vector2, const_b *Vector2) float64 {
	return const_a.X*const_b.X + const_a.Y*const_b.Y
}

func Reflect_Vector2(const_a, const_b *Vector2) Vector2 {
	// make sure n is normalized before using this!
	dot := Dot_Vector2(const_a, const_b)
	return Vector2{
		const_a.X - 2*dot*const_b.X,
		const_a.Y - 2*dot*const_b.Y,
	}
}

type Vector2Int struct {
	X, Y int
}

func DistanceSquare_Vector2Int(const_a, const_b *Vector2) float64 {
	return math.Pow(const_a.X-const_b.X, 2) + math.Pow(const_a.Y-const_b.Y, 2)
}

func Magnitude_Vector2Int(const_a *Vector2Int) float64 {
	return math.Sqrt(float64(const_a.X*const_a.X + const_a.Y*const_a.Y))
}

func MagnitudeSquare_Vector2Int(const_a *Vector2Int) float64 {
	return float64(const_a.X*const_a.X + const_a.Y*const_a.Y)
}

func Normalise_Vector2Int(const_a *Vector2Int) Vector2Int {
	magnitude := Magnitude_Vector2Int(const_a)
	return Vector2Int{int(float64(const_a.X) / magnitude), int(float64(const_a.Y) / magnitude)}
}

func Add_Vector2Int(const_a, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_a.X + const_b.X, const_a.Y + const_b.Y}
}

func Subtract_Vector2Int(const_a, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_a.X - const_b.X, const_a.Y - const_b.Y}
}

func Add_Float_Vector2Int(const_a int, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_b.X + const_a, const_b.Y + const_a}
}

func Multiply_Float_Vector2Int(const_a int, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_b.X * const_a, const_b.Y * const_a}
}
