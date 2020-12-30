package reader

import "github.com/kougazhang/jobx/io"

type IReader interface {
	io.IO
}

type Reader struct {
	IReader
	ReaderSrc io.Src
	ReaderDst io.Dst
}
