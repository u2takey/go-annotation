package args

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
)

// AnnotationArgs is args for annotation-gen tool
type AnnotationArgs struct {
	// tag prefix for an annotation
	AnnotationPrefix string
}


// NewDefaults returns default arguments for the generator.
func NewDefaults() (*args.GeneratorArgs, *AnnotationArgs) {
	genericArgs := args.Default().WithoutDefaultFlagParsing()
	customArgs := &AnnotationArgs{
		AnnotationPrefix: "Annotation@",
	}
	genericArgs.CustomArgs = customArgs // convert to upstream type to make type-casts work there
	genericArgs.OutputFileBaseName = "zz_generated"
	return genericArgs, customArgs
}

// AddFlags add the generator flags to the flag set.
func (ca *AnnotationArgs) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&ca.AnnotationPrefix, "annotation-prefix", ca.AnnotationPrefix, "Tag prefix for an annotation")
}

// Validate checks the given arguments.
func Validate(genericArgs *args.GeneratorArgs) error {
	_ = genericArgs.CustomArgs.(*AnnotationArgs)

	if len(genericArgs.OutputFileBaseName) == 0 {
		return fmt.Errorf("output file base name cannot be empty")
	}

	return nil
}
