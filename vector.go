package main

import (
	"encoding/json"
	"math"
)

type Vector2 struct {
	x, y float64
}

func DistanceSquare_Vector2(const_a, const_b *Vector2) float64 {
	return math.Pow(const_a.x-const_b.x, 2) + math.Pow(const_a.y-const_b.y, 2)
}

func Magnitude_Vector2(const_a *Vector2) float64 {
	return math.Sqrt(const_a.x*const_a.x + const_a.y*const_a.y)
}

func MagnitudeSquare_Vector2(const_a *Vector2) float64 {
	return (const_a.x*const_a.x + const_a.y*const_a.y)
}

func Normalise_Vector2(const_a *Vector2) Vector2 {
	magnitude := Magnitude_Vector2(const_a)
	return Vector2{const_a.x / magnitude, const_a.y / magnitude}
}

func Add_Vector2(const_a, const_b *Vector2) Vector2 {
	return Vector2{const_a.x + const_b.x, const_a.y + const_b.y}
}

func Subtract_Vector2(const_a, const_b *Vector2) Vector2 {
	return Vector2{const_a.x - const_b.x, const_a.y - const_b.y}
}

func Add_Float_Vector2(const_a float64, const_b *Vector2) Vector2 {
	return Vector2{const_b.x + const_a, const_b.y + const_a}
}

func Multiply_Float_Vector2(const_a float64, const_b *Vector2) Vector2 {
	return Vector2{const_b.x * const_a, const_b.y * const_a}
}

func (v *Vector2) UnmarshalJSON(data []byte) error {
	var arr [2]float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.x = arr[0]
	v.y = arr[1]
	return nil
}

type Vector2Int struct {
	x, y int
}

func DistanceSquare_Vector2Int(const_a, const_b *Vector2) float64 {
	return math.Pow(const_a.x-const_b.x, 2) + math.Pow(const_a.y-const_b.y, 2)
}

func Magnitude_Vector2Int(const_a *Vector2Int) float64 {
	return math.Sqrt(float64(const_a.x*const_a.x + const_a.y*const_a.y))
}

func MagnitudeSquare_Vector2Int(const_a *Vector2Int) float64 {
	return float64(const_a.x*const_a.x + const_a.y*const_a.y)
}

func Normalise_Vector2Int(const_a *Vector2Int) Vector2Int {
	magnitude := Magnitude_Vector2Int(const_a)
	return Vector2Int{int(float64(const_a.x) / magnitude), int(float64(const_a.y) / magnitude)}
}

func Add_Vector2Int(const_a, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_a.x + const_b.x, const_a.y + const_b.y}
}

func Subtract_Vector2Int(const_a, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_a.x - const_b.x, const_a.y - const_b.y}
}

func Add_Float_Vector2Int(const_a int, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_b.x + const_a, const_b.y + const_a}
}

func Multiply_Float_Vector2Int(const_a int, const_b *Vector2Int) Vector2Int {
	return Vector2Int{const_b.x * const_a, const_b.y * const_a}
}

func (v *Vector2Int) UnmarshalJSON(data []byte) error {
	var arr [2]int
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.x = arr[0]
	v.y = arr[1]
	return nil
}
