package common

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TableWrite struct {
	o *Option

	table    *tablewriter.Table
	header   []string
	dataRows [][]string
}

func NewTableWrite(w io.Writer, o *Option) *TableWrite {
	t := &TableWrite{
		table:    tablewriter.NewWriter(w),
		dataRows: make([][]string, 0),
		o:        o,
	}

	t.table.SetRowLine(o.TableRowLine)
	t.table.SetAutoMergeCells(o.TableAutoMergeCells)

	return t
}

func (t *TableWrite) SetDataRows(rows [][]string) {
	t.dataRows = rows
}

func (t *TableWrite) SetAutoMergeCells(v bool) {
	t.table.SetAutoMergeCells(v)
}

func (t *TableWrite) SetHeader(header []string) {
	t.header = header
	t.table.SetHeader(header)
}

func (t *TableWrite) SetRowLine(v bool) {
	t.table.SetRowLine(v)
}

func (t *TableWrite) SetAlignCenter(v bool) {
	t.o.TableAlignCenter = v
}

func (t *TableWrite) Display() {
	align := tablewriter.ALIGN_LEFT
	if t.o.TableAlignCenter {
		align = tablewriter.ALIGN_CENTER
	}
	t.table.SetAlignment(align)
	for _, row := range t.dataRows {
		t.table.Append(row)
	}

	if len(t.dataRows) > 0 {
		footer := make([]string, len(t.header))
		footer[len(t.header)-1] = fmt.Sprintf("Total: %d", len(t.dataRows))
		t.table.SetFooter(footer)
	}
	t.table.Render()
}
