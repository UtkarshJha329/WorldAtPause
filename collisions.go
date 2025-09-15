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
	colliderOriginOffset Vector2
	collisionPoints      []Vector2
}

type BoxCollider struct {
	Collider
	size Vector2
}

func (bc *BoxCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	bcTopLeft := Vector2{const_origin.x + bc.colliderOriginOffset.x, const_origin.y + bc.colliderOriginOffset.y}
	bcBottomRight := Vector2{const_origin.x + bc.size.x + bc.colliderOriginOffset.x, const_origin.y + bc.size.y + bc.colliderOriginOffset.y}

	if const_point.x >= bcTopLeft.x && const_point.x <= bcBottomRight.x {
		if const_point.y >= bcTopLeft.y && const_point.y <= bcBottomRight.y {
			return true
		}
	}

	return false
}

func (bc *BoxCollider) CreateCollisionPoints() {

	topLeft := Vector2{bc.colliderOriginOffset.x, bc.colliderOriginOffset.y}
	bottomRight := Vector2{bc.size.x + bc.colliderOriginOffset.x, bc.size.y + bc.colliderOriginOffset.y}
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

func (bc *BoxCollider) GetBoundingBoxDims() *Vector2 {
	return &bc.size
}

func (bc *BoxCollider) GetOffsetOrigin(const_object_origin *Vector2) *Vector2 {
	boxColliderOrigin := Add_Vector2(const_object_origin, &bc.colliderOriginOffset)
	return &boxColliderOrigin
}

type CircleCollider struct {
	Collider
	radius float64
}

func (cc *CircleCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	circleColliderPosition := Add_Vector2(const_origin, &cc.colliderOriginOffset)
	return DistanceSquare_Vector2(const_point, &circleColliderPosition) <= math.Pow(cc.radius, 2.0)
}

func (cc *CircleCollider) CreateCollisionPoints() {
	cc.collisionPoints = []Vector2{cc.colliderOriginOffset}
}

func (cc *CircleCollider) GetCollisionPoints() *[]Vector2 {
	return &cc.collisionPoints
}

func (cc *CircleCollider) CollisionShapeType() int {
	return Circle
}

func (cc *CircleCollider) GetBoundingBoxDims() *Vector2 {
	return &Vector2{cc.radius * 2.0, cc.radius * 2.0}
}

func (cc *CircleCollider) GetOffsetOrigin(const_object_origin *Vector2) *Vector2 {
	circleColliderOrigin := Add_Vector2(const_object_origin, &cc.colliderOriginOffset)
	return &circleColliderOrigin
}

type CircleBoxOverlapCirclePenetrationData struct {
	closestPointOnBox Vector2
	normal            Vector2
	penetrationAmount float64
}

func GetCircleBoxOverlapPenetrationData(const_origin *Vector2, cc *CircleCollider, const_origin_box *Vector2, bx *BoxCollider) CircleBoxOverlapCirclePenetrationData {

	const_offset_origin := cc.GetOffsetOrigin(const_origin)
	closestPointOnBox := *const_offset_origin

	if const_offset_origin.x < const_origin_box.x {
		closestPointOnBox.x = const_origin_box.x
	} else if const_offset_origin.x > const_origin_box.x+bx.size.x {
		closestPointOnBox.x = const_origin_box.x + bx.size.x
	}

	if const_offset_origin.y < const_origin_box.y {
		closestPointOnBox.y = const_origin_box.y
	} else if const_offset_origin.y > const_origin_box.y+bx.size.y {
		closestPointOnBox.y = const_origin_box.y + bx.size.y
	}

	distVector := Vector2{closestPointOnBox.x - const_offset_origin.x, closestPointOnBox.y - const_offset_origin.y}
	return CircleBoxOverlapCirclePenetrationData{
		closestPointOnBox: closestPointOnBox,
		normal:            Normalise_Vector2(&distVector),
		penetrationAmount: cc.radius - Magnitude_Vector2(&distVector),
	}
}

func CircleBoxAABBOverlap(const_origin *Vector2, cc *CircleCollider, const_origin_box *Vector2, bx *BoxCollider) bool {

	const_offset_origin := cc.GetOffsetOrigin(const_origin)
	pointOnBoxToCheck := GetCircleBoxOverlapPenetrationData(const_origin, cc, const_origin_box, bx)
	return DistanceSquare_Vector2(&pointOnBoxToCheck.closestPointOnBox, const_offset_origin) <= math.Pow(cc.radius, 2.0)
}

type CollisionShape interface {
	GetOffsetOrigin(const_object_origin *Vector2) *Vector2
	CollisionShapeType() int
	GetBoundingBoxDims() *Vector2
	GetCollisionPoints() *[]Vector2
	CreateCollisionPoints()
	CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool
}

func CollisionShapeOverlapsWithCollisionShape(const_object_originA *Vector2, csA CollisionShape, const_object_originB *Vector2, csB CollisionShape) bool {

	if csA.CollisionShapeType() == Circle && csB.CollisionShapeType() == Circle {
		csAOffsetOrigin := csA.GetOffsetOrigin(const_object_originA)
		csBOffsetOrigin := csB.GetOffsetOrigin(const_object_originB)

		return DistanceSquare_Vector2(csAOffsetOrigin, csBOffsetOrigin) <= math.Pow(csA.(*CircleCollider).radius+csB.(*CircleCollider).radius, 2.0)
	}
	if csA.CollisionShapeType() == Circle {
		return CircleBoxAABBOverlap(const_object_originA, csA.(*CircleCollider), const_object_originB, csB.(*BoxCollider))
	} else if csB.CollisionShapeType() == Circle {
		return CircleBoxAABBOverlap(const_object_originB, csB.(*CircleCollider), const_object_originA, csA.(*BoxCollider))
	} else {
		return (const_object_originA.x < const_object_originB.x+csB.GetBoundingBoxDims().x) &&
			(const_object_originA.x+csA.GetBoundingBoxDims().x > const_object_originB.x) &&
			(const_object_originA.y < const_object_originB.y+csB.GetBoundingBoxDims().y) &&
			(const_object_originA.y+csA.GetBoundingBoxDims().y > const_object_originB.y)
	}
}

type CollisionShapeAssetData struct {
	CollisionShapeType   string   `json:"CollisionShapeType"`
	Size                 *Vector2 `json:"Size"`
	Radius               *float64 `json:"Radius"`
	ColliderOriginOffset *Vector2 `json:"ColliderOriginOffset"`
}
