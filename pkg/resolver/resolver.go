package resolver

import "glox/pkg/statement"

type Resolver struct {
	pipeline *Pipeline
}

func NewResolver(pipeline *Pipeline) *Resolver {
	return &Resolver{
		pipeline: pipeline,
	}
}

func NewDefaultResolver(table *BindingTable) *Resolver {
	return NewResolverBuilder().WithBinding(table).Build()
}

// Resolver Public Methods
func (r *Resolver) Resolve(statements []statement.Statement) error {
	if r.pipeline == nil {
		return nil
	}
	return r.pipeline.Run(statements)
}

func (r *Resolver) Pipeline() *Pipeline {
	return r.pipeline
}
