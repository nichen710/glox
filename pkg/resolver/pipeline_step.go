package resolver

import "glox/pkg/statement"

type PipelineStep interface {
	Name() string
	Run(statements []statement.Statement) error
}
