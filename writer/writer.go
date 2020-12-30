package writer

import "github.com/kougazhang/jobx/io"

type IWriter interface {
	io.IO
}

type Writer struct {
	IWriter
	WriterSrc io.Src
	WriterDst io.Dst
}
