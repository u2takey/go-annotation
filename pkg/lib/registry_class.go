package lib

import (
	"fmt"
	"reflect"
	"sync"
)

type threadSafeRegistry struct {
	annotationRegistry map[reflect.Type]map[string]Annotation
	mu                 *sync.RWMutex
}

func newThreadSafeRegistry() *threadSafeRegistry {
	return &threadSafeRegistry{
		annotationRegistry: map[reflect.Type]map[string]Annotation{},
		mu:                 &sync.RWMutex{},
	}
}

func (p *threadSafeRegistry) registerAnnotation(t interface{}, a Annotation) {
	annotatedType := reflect.TypeOf(t).Elem()
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.annotationRegistry[annotatedType]; !ok {
		p.annotationRegistry[annotatedType] = map[string]Annotation{}
	}
	annotationMapByType := p.annotationRegistry[annotatedType]

	if _, ok := annotationMapByType[GetAnnotationName(a)]; ok {
		panic(fmt.Sprintf("annotaion plugin with type: %s already registered", annotatedType))
	}
	annotationMapByType[GetAnnotationName(a)] = a
}

func (p *threadSafeRegistry) GetAnnotation(t interface{}, a Annotation) Annotation {
	return p.GetAnnotationByName(t, GetAnnotationName(a))
}

func (p *threadSafeRegistry) GetAnnotationByName(t interface{}, name string) Annotation {
	annotatedType := reflect.TypeOf(t).Elem()
	p.mu.RLock()
	defer p.mu.RUnlock()
	if annotationMapByType, ok := p.annotationRegistry[annotatedType]; ok {
		if annotation, ok := annotationMapByType[name]; ok {
			return annotation
		}
	}
	return nil
}

func (p *threadSafeRegistry) GetAnnotations(t interface{}) (ret []Annotation) {
	annotatedType := reflect.TypeOf(t).Elem()
	p.mu.RLock()
	defer p.mu.RUnlock()
	if annotationMapByType, ok := p.annotationRegistry[annotatedType]; ok {
		for _, v := range annotationMapByType {
			ret = append(ret, v)
		}
	}
	return
}

func (p *threadSafeRegistry) GetAllAnnotations() (ret []Annotation) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for k, annotationMapByType := range p.annotationRegistry {
		if k == reflect.TypeOf(new(Plugin)).Elem() {
			continue
		}
		for _, v2 := range annotationMapByType {
			ret = append(ret, v2)
		}
	}
	return
}

var defaultRegistry = newThreadSafeRegistry()

// RegisterAnnotation method mostly called by generated code to register a annotation with type
func RegisterAnnotation(t interface{}, a Annotation) {
	defaultRegistry.registerAnnotation(t, a)
}

func GetAnnotation(t interface{}, a Annotation) interface{} {
	return defaultRegistry.GetAnnotation(t, a)
}

func GetAnnotations(t interface{}) (ret []Annotation) {
	return defaultRegistry.GetAnnotations(t)
}

func GetAllAnnotations() (ret []Annotation) {
	return defaultRegistry.GetAllAnnotations()
}

type Plugin struct {
}

func RegisterPlugin(a Annotation) {
	defaultRegistry.registerAnnotation(new(Plugin), a)
}

func GetPluginByName(name string) Annotation {
	return defaultRegistry.GetAnnotationByName(new(Plugin), name)
}
