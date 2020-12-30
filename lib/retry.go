package lib

import "time"

type RetryInfo struct {
	Times    int
	Interval time.Duration
}

func Retry(info RetryInfo, fn func() (interface{}, error)) (res interface{}, err error) {
	for i := 0; i < info.Times; i++ {
		res, err = fn()
		if err == nil {
			break
		}
		time.Sleep(info.Interval)
	}
	return res, err
}

func RetryOnlyReturnErr(info RetryInfo, fn func() error) (err error) {
	for i := 0; i < info.Times; i++ {
		err = fn()
		if err == nil {
			break
		}
		time.Sleep(info.Interval)
	}
	return err
}
