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

func (bc *BoxCollider) SetCollisionShapeColliderOriginOffset(offset Vector2) {
	bc.ColliderOriginOffset = offset
}

func (bc *BoxCollider) SetCollisionShapeSize(size Vector2) {
	bc.size = size
}

func (bc *BoxCollider) GetNormalFromPoint(const_object_origin, const_point *Vector2) *Vector2 {
	// boxColliderOrigin := bc.GetOffsetOrigin(const_object_origin)

	topLeft := Add_Vector2(&bc.CollisionPoints[0], const_object_origin)
	bottomRight := Add_Vector2(&bc.CollisionPoints[2], const_object_origin)

	// Distances to each side
	leftDist := const_point.X - topLeft.X
	rightDist := bottomRight.X - const_point.X
	topDist := const_point.Y - topLeft.Y
	bottomDist := bottomRight.Y - const_point.Y

	// Pick the *smallest* distance (least penetration)
	minDist := leftDist
	normal := Vector2{-1, 0}

	if rightDist < minDist {
		minDist = rightDist
		normal = Vector2{1, 0}
	}
	if topDist < minDist {
		minDist = topDist
		normal = Vector2{0, -1}
	}
	if bottomDist < minDist {
		minDist = bottomDist
		normal = Vector2{0, 1}
	}

	return &normal
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

func (cc *CircleCollider) GetNormalFromPoint(const_object_origin *Vector2, const_point *Vector2) *Vector2 {
	circleColliderOrigin := Add_Vector2(const_object_origin, &cc.ColliderOriginOffset)
	normal := Subtract_Vector2(const_point, &circleColliderOrigin)
	return &normal
}

func (cc *CircleCollider) SetCollisionShapeColliderOriginOffset(offset Vector2) {
	cc.ColliderOriginOffset = offset
}

func (cc *CircleCollider) SetCollisionShapeSize(size Vector2) {
	cc.radius = size.X
}

type CollisionShape interface {
	GetOffsetOrigin(const_object_origin *Vector2) *Vector2
	CollisionShapeType() int
	GetBoundingBoxDims() *Vector2
	GetCollisionPoints() *[]Vector2
	GetNormalFromPoint(const_object_origin *Vector2, const_point *Vector2) *Vector2

	SetCollisionShapeSize(size Vector2)
	SetCollisionShapeColliderOriginOffset(offset Vector2)

	CreateCollisionPoints()
	CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool
}

// Collision Checking Functions

func BoxBoxCollisionResult(posA Vector2, boxColA *BoxCollider, posB Vector2, boxColB *BoxCollider) (collisionNormal Vector2, penetrationNormal Vector2, penetration float64, collided bool) {

	boxColACollisionPoints := boxColA.GetCollisionPoints()
	topLeftA := Add_Vector2(&posA, &(*boxColACollisionPoints)[0])
	bottomRightA := Add_Vector2(&posA, &(*boxColACollisionPoints)[2])

	boxColBCollisionPoints := boxColB.GetCollisionPoints()
	topLeftB := Add_Vector2(&posB, &(*boxColBCollisionPoints)[0])
	bottomRightB := Add_Vector2(&posB, &(*boxColBCollisionPoints)[2])

	overlapX := math.Min(bottomRightA.X, bottomRightB.X) - math.Max(topLeftA.X, topLeftB.X)
	overlapY := math.Min(bottomRightA.Y, bottomRightB.Y) - math.Max(topLeftA.Y, topLeftB.Y)

	if overlapX <= 0 || overlapY <= 0 {
		return Vector2{0, 0}, Vector2{0, 0}, 0, false // no collision
	}

	centerA := Vector2{(topLeftA.X + bottomRightA.X) * 0.5, (topLeftA.Y + bottomRightA.Y) * 0.5}
	centerB := Vector2{(topLeftB.X + bottomRightB.X) * 0.5, (topLeftB.Y + bottomRightB.Y) * 0.5}

	if overlapX < overlapY {
		penetration = overlapX
		if centerA.X < centerB.X {
			penetrationNormal = Vector2{-1, 0}
		} else {
			penetrationNormal = Vector2{1, 0}
		}
	} else {
		penetration = overlapY
		if centerA.Y < centerB.Y {
			penetrationNormal = Vector2{0, -1}
		} else {
			penetrationNormal = Vector2{0, 1}
		}
	}

	return penetrationNormal, penetrationNormal, penetration, true
}

func GetBoxCircleOverlapPenetrationData(const_origin_box *Vector2, bx *BoxCollider, const_origin *Vector2, cc *CircleCollider) (collisionNormal Vector2, penetrationNormal Vector2, penetration float64, collided bool) {

	const_offset_origin := cc.GetOffsetOrigin(const_origin)
	closestPointOnBox := *const_offset_origin

	boxCollisionPoints := bx.GetCollisionPoints()
	topLeftBox := Add_Vector2(const_origin_box, &(*boxCollisionPoints)[0])

	if const_offset_origin.X < topLeftBox.X {
		closestPointOnBox.X = topLeftBox.X
	} else if const_offset_origin.X > topLeftBox.X+bx.size.X {
		closestPointOnBox.X = topLeftBox.X + bx.size.X
	}

	if const_offset_origin.Y < topLeftBox.Y {
		closestPointOnBox.Y = topLeftBox.Y
	} else if const_offset_origin.Y > topLeftBox.Y+bx.size.Y {
		closestPointOnBox.Y = topLeftBox.Y + bx.size.Y
	}

	distVector := Vector2{closestPointOnBox.X - const_offset_origin.X, closestPointOnBox.Y - const_offset_origin.Y}
	magDistVector := Magnitude_Vector2(&distVector)
	collisionNormalToReturn := bx.GetNormalFromPoint(const_origin_box, &closestPointOnBox)
	colliding := magDistVector-cc.radius < 0

	return *collisionNormalToReturn, Normalise_Vector2(&distVector), magDistVector - cc.radius, colliding
}

func CollisionShapeOverlapsWithCollisionShape(const_object_originA *Vector2, csA CollisionShape, const_object_originB *Vector2, csB CollisionShape) (collisionNormal Vector2, penetrationNormal Vector2, penetration float64, collided bool) {

	if csA.CollisionShapeType() == Circle && csB.CollisionShapeType() == Circle {
		csAOffsetOrigin := csA.GetOffsetOrigin(const_object_originA)
		csBOffsetOrigin := csB.GetOffsetOrigin(const_object_originB)

		normal := Subtract_Vector2(csAOffsetOrigin, csBOffsetOrigin)
		penetrationAmount := (csA.(*CircleCollider).radius + csB.(*CircleCollider).radius) - Magnitude_Vector2(&normal)
		normal = Normalise_Vector2(&normal)
		collided := penetrationAmount >= 0
		return normal, normal, penetrationAmount, collided
	}
	if csA.CollisionShapeType() == Circle {
		return GetBoxCircleOverlapPenetrationData(const_object_originB, csB.(*BoxCollider), const_object_originA, csA.(*CircleCollider))
	} else if csB.CollisionShapeType() == Circle {
		return GetBoxCircleOverlapPenetrationData(const_object_originA, csA.(*BoxCollider), const_object_originB, csB.(*CircleCollider))
	} else {
		return BoxBoxCollisionResult(*const_object_originA, csA.(*BoxCollider), *const_object_originB, csB.(*BoxCollider))
	}
}

type CollisionShapeAssetData struct {
	CollisionShapeType   string   `json:"CollisionShapeType"`
	Size                 *Vector2 `json:"Size"`
	Radius               *float64 `json:"Radius"`
	ColliderOriginOffset *Vector2 `json:"ColliderOriginOffset"`
}
