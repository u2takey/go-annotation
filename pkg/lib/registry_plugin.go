package lib

type PluginRegistry struct {
	plugins map[string]Annotation
}

// NewPluginRegistry returns PluginRegistry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: map[string]Annotation{},
	}
}

func (r *PluginRegistry) RegisterPlugin(a Annotation) {
	r.plugins[GetAnnotationName(a)] = a
}

func (r *PluginRegistry) GetPlugin(name string) Annotation {
	return r.plugins[name]
}
