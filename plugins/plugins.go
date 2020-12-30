package plugins

import "github.com/kougazhang/jobx/io"

type Plugins struct {
	AfterReaders []func(dst io.Dst) (io.Dst, error)
}
