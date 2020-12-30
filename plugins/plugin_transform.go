package plugins

import (
    "github.com/pkg/errors"
    "github.com/kougazhang/jobx/io"
    "github.com/kougazhang/jobx/lib"
    "time"
)

type ITransform interface {
    Submit(src io.Src, dst io.Dst) (string, error)
    GetStatus(string) (string, error)
    IsFinalStatus(string) (bool, error)
    GetPollingInterval() time.Duration
}

type Transform struct {
    ITransform
    TransformDst io.Dst
    Retry        lib.RetryInfo
}

func (t Transform) SetPluginTransform(src io.Src) func(dst io.Dst) (io.Dst, error) {
    return func(dst io.Dst) (io.Dst, error) {
        taskID, err := t.Submit(src, dst)
        if err != nil {
            return nil, errors.Wrap(err, "ITransform.Submit")
        }

        return dst, t.runSync(taskID)
    }
}

func (t Transform) runSync(taskID string) error {

    for {
        status, err := t.getStatusWithRetry(taskID)
        if err != nil {
            return errors.Wrap(err, "getStatusWithRetry")
        }
        yes, err := t.IsFinalStatus(status)
        if err != nil {
            return errors.Wrap(err, "IsFinalStatus")
        }
        if yes {
            break
        }
        time.Sleep(t.GetPollingInterval())
    }
    return nil
}

func (t Transform) getStatusWithRetry(taskID string) (string, error) {
    res, err := lib.Retry(t.Retry, func() (i interface{}, err error) {
        return t.GetStatus(taskID)
    })
    return res.(string), err
}
