package scale_rules

import "encoding/json"

type ScaleRuler interface {
	Float(value float64) float64
	Int(value int) int
	Parse(rulesString string) (*Rules, error)
	SetRules(rules *Rules)
	SetCombineFunction(combine CombineFunction)
	SetUseOnly(t bool)
}

type CombineFunction int

const (
	CombineLastOne            CombineFunction = 1
	CombineFirstOne           CombineFunction = 2
	CombineRandom             CombineFunction = 3
	CombineAvg                CombineFunction = 4
	CombineMax                CombineFunction = 5
	CombineMin                CombineFunction = 6
	CombineSum                CombineFunction = 7
	CombineIntersections      CombineFunction = 8
	CombineIntersectionPanics CombineFunction = 9
)

type ScaleRules struct {
	rulesContainer *Rules
	combine        CombineFunction
	useResultsOnly *bool
}

func (s *ScaleRules) SetCombineFunction(combine CombineFunction) {
	s.combine = combine
}

func (s *ScaleRules) SetUseOnly(t bool) {
	s.useResultsOnly = &t
}

func (s *ScaleRules) SetRules(rules *Rules) {
	s.rulesContainer = rules
}

func (s *ScaleRules) Float(value float64) float64 {
	return s.rulesContainer.resultForVal(value, s.combine, s.useResultsOnly)
}

func (s *ScaleRules) Int(value int) int {
	return int(s.rulesContainer.resultForVal(float64(value), s.combine, s.useResultsOnly))
}

func (s *ScaleRules) Parse(rulesString string) (*Rules, error) {
	ruls := &Rules{}
	err := json.Unmarshal([]byte(rulesString), ruls)

	return ruls, err
}

func NewScaleRules() *ScaleRules {
	return &ScaleRules{combine: CombineLastOne}
}
