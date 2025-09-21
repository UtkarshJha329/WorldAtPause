package Engine

const (
	QUEST_TYPE_TIME_TRIAL = iota
	QUEST_TYPE_SCORE_LIMIT
	QUEST_TYPE_LEVEL_COMPLETE
)

type QuestData struct {
	QuestType  int
	QuestValue float64
}
