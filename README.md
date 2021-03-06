# Go-annotation
Annotation libraries and tools for golang.


# Docs
Annotation 旨在设计一个适用于 golang 的 annotation 系统，实现类似 Java 的 Annotation 系统，以及常见插件，并提供一定的灵活性用于支持外部插件。

使用 `Annotation@Annotation` 名字 表示使用一个具体的 annotation, 目前内置两个插件 `Description` 和 `Component`

例如, 用 Annotation 系统实现的内置插件 `Component`, 实现了类似 Java 中的依赖注入功能,  具体使用请参考 examples/example_test.go

```golang
// Annotation@Component
type ComponentA struct {
	B1 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B2 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B3 *ComponentB
}

// Annotation@Component={"type": "Singleton"}
type ComponentB struct {
	C *ComponentC `autowired:"true"` // Will populate with NewComponentC()
}

// Annotation@Component
type ComponentC struct {
	D        *ComponentD `autowired:"true"` // Will populate with NewComponentD()
	IntValue int
}

func NewComponentC() *ComponentC {
	return &ComponentC{IntValue: 1}
}

// Annotation@Component
type ComponentD struct {
	IntValue int
}

func NewComponentD() (*ComponentD, error) {
	return &ComponentD{IntValue: 2}, nil
}

```


# RoadMap
[design](/docs/design.md)