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

func (runnerCtx *RunnerContext) HasSharedVariable(fullname string) bool {
	_, exists := runnerCtx.sharedVariables[fullname]
	return exists
}

func (runnerCtx *RunnerContext) ListSharedVariables() []string {
	return slices.Collect(maps.Keys(runnerCtx.sharedVariables))
}

func (runnerCtx *RunnerContext) AddSharedVariable(sharedVars nodes.SharedVariable) {
	fullname := sharedVars.FullName()
	if runnerCtx.HasSharedVariable(fullname) {
		fmt.Printf(`Warning: Overriding shared variable '%s'`, fullname)
	}

	runnerCtx.sharedVariables[fullname] = &sharedVars
}

func (runnerCtx *RunnerContext) GetSharedVariablePtr(fullname string) *nodes.SharedVariable {
	return runnerCtx.sharedVariables[fullname]
}
