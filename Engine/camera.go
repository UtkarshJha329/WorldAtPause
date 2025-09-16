package Engine

type CameraData struct {
	CameraPosition     Vector2
	ScreenSize         Vector2
	TargetFollowOffset Vector2
}

func CameraFollowTarget(targetPosition Vector2, entityComponentsRef *EntityComponents) {

	entityComponentsRef.CameraData.CameraPosition = targetPosition
	currentDirectionToOffsetWorld := Multiply_Float_Vector2(-1.0, &entityComponentsRef.CameraData.CameraPosition)
	halfScreenSize := Multiply_Float_Vector2(0.5, &entityComponentsRef.CameraData.ScreenSize)
	entityComponentsRef.CameraData.TargetFollowOffset = Add_Vector2(&halfScreenSize, &currentDirectionToOffsetWorld)
}
