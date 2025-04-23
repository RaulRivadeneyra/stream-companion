package nodes

type Executor interface {
	Execute(svc *SharedVariableCollection) error
}
