package dbgen

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"sort"
)

type Buffer struct {
	pack             string
	name             string
	goImports        bool
	output           string
	importPackageMap map[string]bool
	structures       []string
}

func (b *Buffer) Convert() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf(`
// Code generated by %s. DO NOT EDIT.

package %s


`, b.name, b.pack))

	if len(b.importPackageMap) > 0 {
		buf.WriteString("import (\n")
		var ps []string
		for p := range b.importPackageMap {
			ps = append(ps, p)
		}
		sort.Slice(ps, func(i, j int) bool {
			return ps[i] < ps[j]
		})
		for _, p := range ps {
			buf.WriteString(fmt.Sprintf("\"%s\"\n", p))
		}
		buf.WriteString(")\n")
	}

	for _, structure := range b.structures {
		buf.WriteString(structure)
	}

	var (
		code []byte
		err  error
	)
	if b.goImports {
		cmd := exec.Command("bash", "-c", "goimports")
		cmd.Stdin = buf
		cmd.Stderr = os.Stderr
		code, err = cmd.Output()
	} else {
		code, err = format.Source(buf.Bytes())
	}

	return code, err
}

type BufferMap map[string]*Buffer

func (m BufferMap) get(output string, opts Options) *Buffer {
	buf := m[output]
	if buf == nil {
		buf = &Buffer{
			output:           output,
			importPackageMap: make(map[string]bool),
			pack:             opts.GoPackage,
			name:             opts.GeneratorName,
			goImports:        opts.EnableGoImports,
		}
		m[output] = buf
	}

	return buf
}
