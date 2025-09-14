package main

type CameraData struct {
	cameraPosition     Vector2
	screenSize         Vector2
	targetFollowOffset Vector2
}

func CameraFollowTarget(targetPosition Vector2, entityComponentsRef *EntityComponents) {

	entityComponentsRef.cameraData.cameraPosition = targetPosition
	currentDirectionToOffsetWorld := Multiply_Float_Vector2(-1.0, &entityComponentsRef.cameraData.cameraPosition)
	halfScreenSize := Multiply_Float_Vector2(0.5, &entityComponentsRef.cameraData.screenSize)
	entityComponentsRef.cameraData.targetFollowOffset = Add_Vector2(&halfScreenSize, &currentDirectionToOffsetWorld)
}
