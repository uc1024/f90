package main

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
)

func genProto(gen *protogen.Plugin) error {
	gen.SupportedFeatures = uint64(PLUGINPB_FEATURE)
	lo.ForEach(gen.Files, func(item *protogen.File, i int) {
		if item.Generate {
			generateFile(gen, item, false)
		}
	})
	return nil
}

// generateFile create a .ginx.pb.go file.
func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + FILE_SUFFIX
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	return g
}
