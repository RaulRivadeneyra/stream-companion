package branch

import "github.com/RaulRivadeneyra/stream-companion/internal/nodes"

type IfBranch struct {
}

type CaseBranch struct {
}

type Branch interface {
	IfBranch | CaseBranch
}

type BranchBuilder struct {
	branchType BranchType
	nextNode   *nodes.ExecutionNode
}

func NewBranchBuilder() *BranchBuilder {
	return &BranchBuilder{}
}
