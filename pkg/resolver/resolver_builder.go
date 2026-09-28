package resolver

type ResolverBuilder struct {
	bindingTable *BindingTable
	customSteps  []PipelineStep
}

func NewResolverBuilder() *ResolverBuilder {
	return &ResolverBuilder{
		customSteps: make([]PipelineStep, 0),
	}
}

func (b *ResolverBuilder) WithBinding(table *BindingTable) *ResolverBuilder {
	b.bindingTable = table
	return b
}

func (b *ResolverBuilder) AddStep(step PipelineStep) *ResolverBuilder {
	if step != nil {
		b.customSteps = append(b.customSteps, step)
	}
	return b
}

func (b *ResolverBuilder) Build() *Resolver {
	var steps []PipelineStep
	if b.bindingTable != nil {
		steps = append(steps, NewBindingStep(b.bindingTable))
	}
	steps = append(steps, b.customSteps...)

	return NewResolver(NewPipeline(steps...))
}
