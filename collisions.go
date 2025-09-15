package main

import (
	"math"
)

type ColliderType int

const (
	Circle = iota
	Box
)

type Collider struct {
	collisionPoints []Vector2
}

type BoxCollider struct {
	Collider
	size Vector2
}

func (bc *BoxCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	bcTopLeft := Vector2{const_origin.x, const_origin.y}
	bcBottomRight := Vector2{const_origin.x + bc.size.x, const_origin.y + bc.size.y}

	if const_point.x >= bcTopLeft.x && const_point.x <= bcBottomRight.x {
		if const_point.y >= bcTopLeft.y && const_point.y <= bcBottomRight.y {
			return true
		}
	}

	return false
}

func (bc *BoxCollider) CreateCollisionPoints() {

	topLeft := Vector2{0, 0}
	bottomRight := bc.size
	topRight := Vector2{bottomRight.x, topLeft.y}
	bottomLeft := Vector2{topLeft.x, bottomRight.y}

	bc.collisionPoints = []Vector2{
		topLeft,
		topRight,
		bottomRight,
		bottomLeft,
	}
}

func (bc *BoxCollider) GetCollisionPoints() *[]Vector2 {
	return &bc.collisionPoints
}

func (bc *BoxCollider) CollisionShapeType() int {
	return Box
}

type CircleCollider struct {
	Collider
	radius float64
}

func (cc *CircleCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	return DistanceSquare_Vector2(const_point, const_origin) <= math.Pow(cc.radius, 2.0)
}

func (cc *CircleCollider) CreateCollisionPoints() {
	cc.collisionPoints = []Vector2{{0, 0}}
}

func (cc *CircleCollider) GetCollisionPoints() *[]Vector2 {
	return &cc.collisionPoints
}

func (cc *CircleCollider) CollisionShapeType() int {
	return Circle
}

func CircleBoxAABBOverlap(const_origin_circle *Vector2, cc *CircleCollider, const_origin_box *Vector2, bx *BoxCollider) bool {

	pointOnBoxToCheck := *const_origin_circle

	if const_origin_circle.x < const_origin_box.x {
		pointOnBoxToCheck.x = const_origin_box.x
	} else if const_origin_circle.x > const_origin_box.x+bx.size.x {
		pointOnBoxToCheck.x = const_origin_box.x + bx.size.x
	}

	if const_origin_circle.y < const_origin_box.y {
		pointOnBoxToCheck.y = const_origin_box.y
	} else if const_origin_circle.y > const_origin_box.y+bx.size.y {
		pointOnBoxToCheck.y = const_origin_box.y + bx.size.y
	}

	return DistanceSquare_Vector2(&pointOnBoxToCheck, const_origin_circle) <= math.Pow(cc.radius, 2.0)
}

type CollisionShape interface {
	CollisionShapeType() int
	GetCollisionPoints() *[]Vector2
	CreateCollisionPoints()
	CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool
}

func CollisionShapeOverlapsWithCollisionShape(const_originA *Vector2, csA CollisionShape, const_originB *Vector2, csB CollisionShape) bool {

	if csA.CollisionShapeType() == Circle && csB.CollisionShapeType() == Circle {
		return DistanceSquare_Vector2(const_originA, const_originB) <= math.Pow(csA.(*CircleCollider).radius+csB.(*CircleCollider).radius, 2.0)
	}
	if csA.CollisionShapeType() == Circle {
		return CircleBoxAABBOverlap(const_originA, csA.(*CircleCollider), const_originB, csB.(*BoxCollider))
	} else if csB.CollisionShapeType() == Circle {
		return CircleBoxAABBOverlap(const_originB, csB.(*CircleCollider), const_originA, csA.(*BoxCollider))
	} else {
		return (const_originA.x < const_originB.x+csB.(*BoxCollider).size.x) &&
			(const_originA.x+csA.(*BoxCollider).size.x > const_originB.x) &&
			(const_originA.y < const_originB.y+csB.(*BoxCollider).size.y) &&
			(const_originA.y+csA.(*BoxCollider).size.y > const_originB.y)
	}
}

type CollisionShapeAssetData struct {
	CollisionShapeType string   `json:"CollisionShapeType"`
	Size               *Vector2 `json:"Size"`
	Radius             *float64 `json:"Radius"`
}
