package Engine

import "time"

type GameStatsData struct {
	TimeSinceLaunch time.Duration
	TimeLastFrame   time.Duration
}
