# Design

## Annotation

格式: `prefix$Type=$Body`

- prefix: 前缀 比如 Annotation@
- Type: annotation 类型
- Body: annotation 内容，使用 json 格式

```go
type MyAnnotation struct{
    A string
}


// Annotation@MyAnnotation={"a": 1}
type A struct{}
```

## RoadMap

- lib: 支持常见的函数，比如:
    - [ ] RegisterAnnotation(classType, annotation)
    - [ ] GetAnnotation(classType, annotationType) annotation
    - [ ] GetAnnotations(classType) []annotation
    - [ ] GetAllAnnotations() []annotation
- 内置 annotation：提供生成模版 和一些检查 过滤函数
- code generator tool，生成 注册函数
    - [ ] 检查格式
    - [ ] 调用 RegisterAnnotation
    - [ ] 生成 (TypeA)GetAnnotationWithTypeA() annotationType
    - [ ] 处理 内置 annotation，完成内置 annotation 的代码生成

实现: 对照 Java

- Retention:
    - [ ] Runtime
- Target:
    - [ ] ElementType.TYPE
    - [ ] ElementType.FIELD
    - [ ] ElementType.METHOD
    - [ ] ElementType.PARAMETER
    - [ ] ElementType.LOCAL_VARIABLE
    - [ ] ElementType.PACKAGE
        
    


