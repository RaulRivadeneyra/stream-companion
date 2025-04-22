package runner

import (
	"fmt"
	"maps"
	"slices"

	"github.com/RaulRivadeneyra/stream-companion/internal/nodes"
)

type RunnerContext struct {
	sharedVariables map[string]*nodes.SharedVariable
}

func (rc *RunnerContext) HasSharedVariable(fullname string) bool {
	_, exists := rc.sharedVariables[fullname]
	return exists
}

func (rc *RunnerContext) ListSharedVariables() []string {
	return slices.Collect(maps.Keys(rc.sharedVariables))
}

func (rc *RunnerContext) AddSharedVariable(sv nodes.SharedVariable) {
	fullname := sv.FullName()
	if rc.HasSharedVariable(fullname) {
		fmt.Printf(`Warning: Overriding shared variable '%s'`, fullname)
	}

	rc.sharedVariables[fullname] = &sv
}

func (rc *RunnerContext) GetSharedVariablePtr(fullname string) *nodes.SharedVariable {
	return rc.sharedVariables[fullname]
}
