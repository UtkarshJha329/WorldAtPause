package main

type CameraData struct {
	screenSize         Vector2
	targetFollowOffset Vector2
}

func CameraFollowTarget(targetPosition Vector2, entityComponentsRef *EntityComponents) {

	entityComponentsRef.positions[entityComponentsRef.cameraEntityID] = targetPosition
	currentDirectionToOffsetWorld := Multiply_Float_Vector2(-1.0, &entityComponentsRef.positions[entityComponentsRef.cameraEntityID])
	halfScreenSize := Multiply_Float_Vector2(0.5, &entityComponentsRef.cameraData.screenSize)
	entityComponentsRef.cameraData.targetFollowOffset = Add_Vector2(&halfScreenSize, &currentDirectionToOffsetWorld)
}
