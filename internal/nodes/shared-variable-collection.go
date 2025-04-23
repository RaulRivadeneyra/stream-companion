package nodes

import (
	"fmt"
	"maps"
	"slices"
)

type SharedVariableCollection struct {
	data map[string]*SharedVariable
}

func NewSharedVariableCollection() *SharedVariableCollection {
	svc := SharedVariableCollection{}
	svc.data = make(map[string]*SharedVariable)
	return &svc
}

func (svc *SharedVariableCollection) HasSharedVariable(fullname string) bool {
	_, exists := svc.data[fullname]
	return exists
}

func (svc *SharedVariableCollection) ListSharedVariables() []string {
	return slices.Collect(maps.Keys(svc.data))
}

func (svc *SharedVariableCollection) AddSharedVariable(sv *SharedVariable) {
	fullname := sv.FullName()
	if svc.HasSharedVariable(fullname) {
		fmt.Printf(`Warning: Overriding shared variable '%s'`, fullname)
	}

	svc.data[fullname] = sv
}

func (svc *SharedVariableCollection) GetSharedVariable(fullname string) SharedVariable {
	return *svc.data[fullname]
}

func (svc *SharedVariableCollection) GetSharedVariableList() []*SharedVariable {
	return slices.Collect(maps.Values(svc.data))
}
