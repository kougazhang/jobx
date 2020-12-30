package job

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/kougazhang/jobx/trigger"
)

const (
	TriggerRelationAnd = "and"
	TriggerRelationOr  = "or"
)

type Trigger struct {
	Conditions []trigger.Trigger
	Relation   string
}

func (t Trigger) IsValidRelation() bool {
	return t.Relation == TriggerRelationAnd || t.Relation == TriggerRelationOr
}

func (t Trigger) paincInvalidRelation() {
	panic(fmt.Sprintf("trigger relation %s invalid", t.Relation))
}

func (t Trigger) Trigger() (can bool, err error) {
	if !t.IsValidRelation() {
		t.paincInvalidRelation()
	}

	for _, conditionFn := range t.Conditions {
		can, err = conditionFn.CanStart()
		if err != nil {
			return can, errors.Wrap(err, "Filter")
		}
		if t.Relation == TriggerRelationAnd {
			if !can {
				return can, nil
			}
		} else if t.Relation == TriggerRelationOr {
			if can {
				return can, nil
			}
		} else {
			t.paincInvalidRelation()
		}
	}

	return can, err
}
