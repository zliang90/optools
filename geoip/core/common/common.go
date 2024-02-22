package common

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

type Option struct {
	ShowMyExternal bool
	MyExternalAddr string
	IpSets         sets.String
	OutToJSON      bool

	// for table writer
	TableRowLine        bool
	TableAutoMergeCells bool
	TableAlignCenter    bool
}

func DefaultOption() *Option {
	return &Option{
		ShowMyExternal: true,
		IpSets:         sets.NewString(),
		OutToJSON:      true,

		TableRowLine:        true,
		TableAutoMergeCells: true,
		TableAlignCenter:    false,
	}
}
