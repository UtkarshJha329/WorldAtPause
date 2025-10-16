package Engine

func LinearInterpolationInt(interpolationParameter *float64, startValue int, curValueToChange *int, finalValue *int) {
	newValue := (startValue) + int(float64(*finalValue-startValue)*(*interpolationParameter))
	*curValueToChange = newValue
}

func LinearInterpolationFloat(interpolationParameter *float64, startValue float64, curValueToChange *float64, finalValue *float64) {
	newValue := (startValue) + (*finalValue-startValue)*(*interpolationParameter)
	*curValueToChange = newValue
}

func LinearInterpolationVector2(interpolationParameter *float64, startValue Vector2, curValueToChange *Vector2, finalValue *Vector2) {
	totalDeltaFromStartToFinish := Subtract_Vector2(finalValue, &startValue)
	scaledDelta := Multiply_Float_Vector2(*interpolationParameter, &totalDeltaFromStartToFinish)
	newValue := Add_Vector2(&startValue, &scaledDelta)
	*curValueToChange = newValue
}

func LinearInterpolationVector2Int(interpolationParameter *float64, startValue Vector2Int, curValueToChange *Vector2Int, finalValue *Vector2Int) {
	totalDeltaFromStartToFinish := Subtract_Vector2Int(finalValue, &startValue)
	scaledDelta := Multiply_Float_Vector2Int(*interpolationParameter, &totalDeltaFromStartToFinish)
	newValue := Add_Vector2Int(&startValue, &scaledDelta)
	*curValueToChange = newValue
}
