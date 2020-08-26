package scale_rules

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Rules []Rule

type Rule string

func (r *Rules) resultForVal(value float64, combine CombineFunction, only *bool) float64 {
	var results []float64
	var minResult float64
	var maxResult float64
	var sumResult float64
	var intersectionFlag bool
	var intersectionUsed int

	for _, rule := range *r {
		result, err, used := rule.result(value)
		if err != nil {
			log.Print("error in scale rule lib: ", err)
		}

		if only != nil && *only == false && used {
			continue
		}

		if only != nil && *only == true && !used {
			continue
		}

		if used && intersectionFlag {
			intersectionUsed++
		} else if used && !intersectionFlag {
			intersectionFlag = true
		}

		if minResult > result {
			minResult = result
		}

		if maxResult < result {
			maxResult = result
		}

		sumResult += result

		results = append(results, result)
	}

	if results == nil || len(results) == 0 {
		return 0
	}

	switch combine {
	case CombineLastOne:
		return results[len(results)-1]
	case CombineFirstOne:
		return results[0]
	case CombineRandom:
		rand.Seed(time.Now().UnixNano())
		return results[rand.Intn(len(results)-1)]
	case CombineAvg:
		return sumResult / float64(len(results))
	case CombineMax:
		return maxResult
	case CombineMin:
		return minResult
	case CombineSum:
		return sumResult
	case CombineIntersections:
		return float64(intersectionUsed)
	case CombineIntersectionPanics:
		if intersectionUsed > 0 {
			panic("scale rule panics")
		}
	}

	return 0
}

const (
	eqlPrefix = "="
	gtrPrefix = ">"
	smrPrefix = "<"
	notPrefix = "!"
	rngPrefix = "@"
)

func (r *Rule) result(value float64) (float64, error, bool) {
	var expression = regexp.MustCompile(`^(?P<Manager>\W)(?P<Condition>.*)\?(?P<True>.*):(?P<False>.*)$`)
	var err error
	var condition float64

	match := expression.FindStringSubmatch(string(*r))

	if len(match) != 5 {
		return 0, errors.New("expression not parsable (" + string(*r) + ")"), false
	}

	if match[1] != rngPrefix {
		condition, err = strconv.ParseFloat(match[2], 64)
		if err != nil {
			return 0, err, false
		}
	}

	trueVal, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		return 0, err, false
	}
	falseVal, err := strconv.ParseFloat(match[4], 64)
	if err != nil {
		return 0, err, false
	}

	return handleCondition(value, match, condition, trueVal, falseVal)
}

func handleCondition(value float64, match []string, condition float64, trueVal float64, falseVal float64) (float64, error, bool) {
	switch match[1] {
	case smrPrefix:
		if value < condition {
			return trueVal, nil, true
		}

		return falseVal, nil, false
	case gtrPrefix:
		if value > condition {
			return trueVal, nil, true
		}

		return falseVal, nil, false
	case eqlPrefix:
		if value == condition {
			return trueVal, nil, true
		}

		return falseVal, nil, false
	case notPrefix:
		if value != condition {
			return trueVal, nil, true
		}

		return falseVal, nil, false
	case rngPrefix:
		conditions := strings.Split(match[2], "-")
		if len(conditions) == 2 {
			conditionFrom, err := strconv.ParseFloat(conditions[0], 64)
			if err != nil {
				return 0, err, true
			}
			conditionTo, err := strconv.ParseFloat(conditions[0], 64)
			if err != nil {
				return 0, err, true
			}

			if value >= conditionFrom && value <= conditionTo {
				return trueVal, nil, true
			}

			return falseVal, nil, false

		} else {
			return 0, errors.New("between conditions error"), true
		}
	}
	return 0, nil, false
}
