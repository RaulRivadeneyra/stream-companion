package runner

import "github.com/RaulRivadeneyra/stream-companion/internal/nodes"

type Variables = map[string]any
type Condition = func(input Variables) (bool, error)

type Runner struct {
	sharedVariables Variables
	startingNode    nodes.INode[any]
}

func (r *Runner) SetStartingNode(node nodes.INode[any]) {
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
