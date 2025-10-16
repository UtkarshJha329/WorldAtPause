package CrossyRoad

import (
	"math/rand/v2"
)

const (
	CrossyRoads_SectionType_ROAD = iota
	CrossyRoads_SectionType_POOL
	CrossyRoads_SectionType_PARKING
	CrossyRoads_SectionType_SPACE
	CrossyRoads_SectionType_FINALE
)

type Section struct {
	SectionType                  int
	SectionForNumRows            int
	SectionScrolling             bool
	NumSectionRowsWrittenToBoard int
}

// MINIMUM GENERATED SECTION SIZE IS AT LEAST 5
func (section *Section) GenerateRandomSection() {
	section.SectionType = rand.IntN(5)
	section.SectionForNumRows = 5 + rand.IntN(5)
	section.SectionScrolling = rand.IntN(2) == 1
}

func (section *Section) GenerateRandomSectionWithNumRows(numRows int) {
	section.SectionType = rand.IntN(5)
	section.SectionForNumRows = numRows
	section.SectionScrolling = rand.IntN(2) == 1
}

func (section *Section) GenerateRandomSectionWithNumRowsAndOfType(sectionType int, numRows int) {
	section.SectionType = sectionType
	section.SectionForNumRows = numRows
	section.SectionScrolling = rand.IntN(2) == 1
}

func (section *Section) GenerateSectionWithNumRowsOfTypeAndLoopType(sectionType int, numRows int, sectionLooping bool) {
	section.SectionType = sectionType
	section.SectionForNumRows = numRows
	section.SectionScrolling = sectionLooping
}

func (section *Section) GenerateRandomSectionButNotSection(minimumNumberOfRowsInSection int, additionalMaxRandomNumberOfRowsInSection int, sectionTypeToAvoid int) {
	section.SectionType = rand.IntN(CrossyRoads_SectionType_FINALE)
	for section.SectionType == sectionTypeToAvoid {
		section.SectionType = rand.IntN(CrossyRoads_SectionType_FINALE)
	}
	section.SectionForNumRows = minimumNumberOfRowsInSection + rand.IntN(additionalMaxRandomNumberOfRowsInSection)
	section.SectionScrolling = rand.IntN(2) == 1

	// fmt.Println("Creating section of type := ", section.SectionType)
}
