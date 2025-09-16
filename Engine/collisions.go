package Engine

import (
	"math"
)

type ColliderType int

const (
	Circle = iota
	Box
)

type Collider struct {
	ColliderOriginOffset Vector2
	CollisionPoints      []Vector2
}

type BoxCollider struct {
	Collider
	size Vector2
}

func (bc *BoxCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	bcTopLeft := Vector2{const_origin.X + bc.ColliderOriginOffset.X, const_origin.Y + bc.ColliderOriginOffset.Y}
	bcBottomRight := Vector2{const_origin.X + bc.size.X + bc.ColliderOriginOffset.X, const_origin.Y + bc.size.Y + bc.ColliderOriginOffset.Y}

	if const_point.X >= bcTopLeft.X && const_point.X <= bcBottomRight.X {
		if const_point.Y >= bcTopLeft.Y && const_point.Y <= bcBottomRight.Y {
			return true
		}
	}

	return false
}

func (bc *BoxCollider) CreateCollisionPoints() {

	topLeft := Vector2{bc.ColliderOriginOffset.X, bc.ColliderOriginOffset.Y}
	bottomRight := Vector2{bc.size.X + bc.ColliderOriginOffset.X, bc.size.Y + bc.ColliderOriginOffset.Y}
	topRight := Vector2{bottomRight.X, topLeft.Y}
	bottomLeft := Vector2{topLeft.X, bottomRight.Y}

	bc.CollisionPoints = []Vector2{
		topLeft,
		topRight,
		bottomRight,
		bottomLeft,
	}
}

func (bc *BoxCollider) GetCollisionPoints() *[]Vector2 {
	return &bc.CollisionPoints
}

func (bc *BoxCollider) CollisionShapeType() int {
	return Box
}

func (bc *BoxCollider) GetBoundingBoxDims() *Vector2 {
	return &bc.size
}

func (bc *BoxCollider) GetOffsetOrigin(const_object_origin *Vector2) *Vector2 {
	boxColliderOrigin := Add_Vector2(const_object_origin, &bc.ColliderOriginOffset)
	return &boxColliderOrigin
}

type CircleCollider struct {
	Collider
	radius float64
}

func (cc *CircleCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	circleColliderPosition := Add_Vector2(const_origin, &cc.ColliderOriginOffset)
	return DistanceSquare_Vector2(const_point, &circleColliderPosition) <= math.Pow(cc.radius, 2.0)
}

func (cc *CircleCollider) CreateCollisionPoints() {
	cc.CollisionPoints = []Vector2{cc.ColliderOriginOffset}
}

func (cc *CircleCollider) GetCollisionPoints() *[]Vector2 {
	return &cc.CollisionPoints
}

func (cc *CircleCollider) CollisionShapeType() int {
	return Circle
}

func (cc *CircleCollider) GetBoundingBoxDims() *Vector2 {
	return &Vector2{cc.radius * 2.0, cc.radius * 2.0}
}

func (cc *CircleCollider) GetOffsetOrigin(const_object_origin *Vector2) *Vector2 {
	circleColliderOrigin := Add_Vector2(const_object_origin, &cc.ColliderOriginOffset)
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

	if const_offset_origin.X < const_origin_box.X {
		closestPointOnBox.X = const_origin_box.X
	} else if const_offset_origin.X > const_origin_box.X+bx.size.X {
		closestPointOnBox.X = const_origin_box.X + bx.size.X
	}

	if const_offset_origin.Y < const_origin_box.Y {
		closestPointOnBox.Y = const_origin_box.Y
	} else if const_offset_origin.Y > const_origin_box.Y+bx.size.Y {
		closestPointOnBox.Y = const_origin_box.Y + bx.size.Y
	}

	distVector := Vector2{closestPointOnBox.X - const_offset_origin.X, closestPointOnBox.Y - const_offset_origin.Y}
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
		return (const_object_originA.X < const_object_originB.X+csB.GetBoundingBoxDims().X) &&
			(const_object_originA.X+csA.GetBoundingBoxDims().X > const_object_originB.X) &&
			(const_object_originA.Y < const_object_originB.Y+csB.GetBoundingBoxDims().Y) &&
			(const_object_originA.Y+csA.GetBoundingBoxDims().Y > const_object_originB.Y)
	}
}

type CollisionShapeAssetData struct {
	CollisionShapeType   string   `json:"CollisionShapeType"`
	Size                 *Vector2 `json:"Size"`
	Radius               *float64 `json:"Radius"`
	ColliderOriginOffset *Vector2 `json:"ColliderOriginOffset"`
}
