package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	annoArgs "github.com/u2takey/go-annotation/cmd/annotation-gen/args"
	"github.com/u2takey/go-annotation/pkg/lib"
	_ "github.com/u2takey/go-annotation/pkg/plugins"
)

// prefix$Enable=true
// prefix$Type=$Body
type annotation struct {
	rawTypeName string
	body        string
}

// key is annotation rawTypeName
type annotations map[string]*annotation

func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func annotationEnabled(prefix string, comments []string) bool {
	enableFlag := prefix + "Enable"
	for _, comment := range comments {
		if strings.HasPrefix(comment, enableFlag) {
			return true
		}
	}
	return false
}

func extractAnnotations(prefix string, t *types.Type) annotations {
	ret := map[string]*annotation{}
	comments := append(t.CommentLines, t.SecondClosestCommentLines...)
	for _, comment := range comments {
		var rawTypeName, body string
		if !strings.HasPrefix(comment, prefix) {
			continue
		}
		s := strings.TrimPrefix(comment, prefix)
		sl := strings.Split(s, "=")
		if len(sl) == 1 {
			rawTypeName, body = sl[0], ""
		} else if len(sl) == 2 {
			rawTypeName, body = sl[0], sl[1]
		} else {
			klog.V(4).Infof("annotation format not valid %s\n", comment)
			continue
		}

		if !isJson(body) {
			klog.V(1).Infoln("annotation format not valid: not valid json", body)
			continue
		}

		ret[rawTypeName] = &annotation{
			rawTypeName: rawTypeName,
			body:        body,
		}
	}
	return ret
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPrivateNamer(0, ""),
		"raw":    namer.NewRawNamer("", nil),
	}
}

func DefaultNameSystem() string {
	return "public"
}

// Packages
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	// LoadGoBoilerplate
	//boilerplate, err := arguments.LoadGoBoilerplate()
	//if err != nil {
	//	klog.Fatalf("Failed loading boilerplate: %v", err)
	//}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	annotationArgs := arguments.CustomArgs.(*annoArgs.AnnotationArgs)

	// header
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)) /*boilerplate...*/)

	// arguments handling

	// inputs, get package from context.Universe
	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]

		for _, a := range pkg.Imports {
			context.AddDirectory(a.Path)
		}
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}
		//
		if !annotationEnabled(annotationArgs.AnnotationPrefix, pkg.Comments) {
			continue
		}

		klog.V(5).Infof("Generating for pkg %q", i)

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				HeaderText:  header,
				// generator 一个 Generator 生成一个文件
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						NewGenAnnotation(arguments.OutputFileBaseName, pkg.Path, annotationArgs.AnnotationPrefix),
					}
				},
				// 过滤函数，哪些 type 不关心的，直接过滤，不会被 generator 处理
				// generator 的 过滤器 也可以完成类似的事情，调用时机不同
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})

	}
	return packages
}

// Order
// 1. Filter()        // Subsequent calls see only types that pass this.
// 2. Namers()        // Subsequent calls see the namers provided by this.
// 3. PackageVars()
// 4. PackageConsts()
// 5. Init()
// 6. GenerateType()  // Called N times, once per type in the context's Order.
// 7. Imports()
type genAnnotation struct {
	generator.DefaultGen
	targetPackage    string
	annotationPrefix string
	imports          namer.ImportTracker
	typesForInit     []*types.Type
}

func NewGenAnnotation(sanitizedName, targetPackage, annotationPrefix string) generator.Generator {
	return &genAnnotation{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		targetPackage:    targetPackage,
		annotationPrefix: annotationPrefix,
		imports:          generator.NewImportTracker(),
		typesForInit:     make([]*types.Type, 0),
	}
}

// Namer for template
func (g *genAnnotation) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

//
func (g *genAnnotation) Filter(c *generator.Context, t *types.Type) bool {
	return true
}

//
func (g *genAnnotation) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "github.com/u2takey/go-annotation/pkg/lib")
	return
}

// Init method for generated code
func (g *genAnnotation) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func findAnnotationType(c *generator.Context, name string) *types.Type {
	for _, p := range c.Universe {
		for _, t := range p.Types {
			klog.V(8).Infoln("finding name", t.Name.Name)
			if t.Name.Name == name {
				return t
			}
		}
	}
	return nil
}

// core
func (g *genAnnotation) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	// writer
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.V(5).Infof("processing type %v", t)
	// params
	annotations := extractAnnotations(g.annotationPrefix, t)
	for _, anno := range annotations {
		annotationType := findAnnotationType(c, anno.rawTypeName)
		if annotationType == nil {
			klog.V(1).Infoln("annotation type not found", anno.rawTypeName)
			continue
		}
		m := map[string]interface{}{
			"Resource":       c.Universe.Function(types.Name{Package: t.Name.Package, Name: "Resource"}),
			"type":           t,
			"annotationType": annotationType,
			"annotationBody": anno.body,
		}
		klog.V(3).Infoln("annotation m", m)
		// render
		sw.Do(typeListerInterface, m)
		annotationPlugin := lib.GetPluginByName(anno.rawTypeName)

		if compile, ok := annotationPlugin.(lib.CompileAnnotation); ok {
			sw.Do(compile.GetTemplate(), m)
		} else {
			klog.V(4).Infoln("get compile", annotationPlugin)
		}
	}

	return sw.Error()
}

// register template
var typeListerInterface = `
func init() {
	b := new($.annotationType|raw$)
	err := json.Unmarshal([]byte("$.annotationBody|js$"), b)
	if err != nil {
		klog.Fatal("unmarshal json failed", err)
		return
	}
	lib.RegisterAnnotation(new($.type|raw$), b)
}

`
