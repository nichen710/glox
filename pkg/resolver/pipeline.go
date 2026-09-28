package resolver

import "glox/pkg/statement"

type Pipeline struct {
	steps []PipelineStep
}

func NewPipeline(steps ...PipelineStep) *Pipeline {
	return &Pipeline{steps: steps}
}

func (p *Pipeline) Run(statements []statement.Statement) error {
	for _, step := range p.steps {
		if err := step.Run(statements); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) Steps() []PipelineStep {
	return p.steps
}
