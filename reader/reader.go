package reader

import "jobx/io"

type IReader interface {
	io.IO
}

type Reader struct {
	IReader
	ReaderSrc io.Src
	ReaderDst io.Dst
}
