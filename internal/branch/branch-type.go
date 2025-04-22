package branch

type BranchType int

const (
	IfStatement BranchType = iota
	CaseStatement
)
