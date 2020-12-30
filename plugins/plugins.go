package plugins

import "jobx/io"

type Plugins struct {
	AfterReaders []func(dst io.Dst) (io.Dst, error)
}
