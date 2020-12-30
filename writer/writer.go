package writer

import "jobx/io"

type IWriter interface {
	io.IO
}

type Writer struct {
	IWriter
	WriterSrc io.Src
	WriterDst io.Dst
}
