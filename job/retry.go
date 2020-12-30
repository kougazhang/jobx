package job

import "jobx/lib"

type Retry struct {
	Trigger            lib.RetryInfo
	Job                lib.RetryInfo
	GetTransformStatus lib.RetryInfo
}
