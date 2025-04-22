package runner

import "github.com/RaulRivadeneyra/stream-companion/internal/nodes"

type Variables = map[string]any
type Condition = func(input Variables) (bool, error)

type BranchingNode struct {
	branches map[any]nodes.ExecutionNode
}

type Runner struct {
	sharedVariables Variables
	startingNode    *nodes.ExecutionNode
}

func (r *Runner) SetStartingNode(node *nodes.ExecutionNode) {
	r.startingNode = node
}

const MAX_NODES_EXECUTED = 1000

type RunnerResult struct {
	status string

	error error
}

func (r *Runner) Run() (any, error) {
	// currentNode := r.startingNode
	// nodesExecuted := 0
	// testing := uuid.New()
	// for {
	// 	if nodesExecuted >= MAX_NODES_EXECUTED {
	// 		break
	// 	}
	// }

	return nil, nil
}

type verificationResult = struct {
	isValid bool
	issues  []string
}

func (r *Runner) Verify() verificationResult {
	return verificationResult{
		isValid: true,
	}
}
