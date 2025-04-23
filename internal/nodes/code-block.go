package nodes

import "os"

type CodeBlock struct {
	body string
}

func (cb *CodeBlock) BodyFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cb.body = string(content)
	return nil
}
