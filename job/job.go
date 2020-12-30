package job

import (
    "github.com/kougazhang/jobx/hook"
    "github.com/kougazhang/jobx/lib"
    "github.com/kougazhang/jobx/reader"
    "github.com/kougazhang/jobx/writer"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "time"
)

type Job struct {
    Log     *log.Entry
    TraceID string

    // trigger Job
    TriggerConditions *Trigger

    // job read
    reader.Reader

    // job write
    Writer *writer.Writer

    // retry for io
    IORetry Retry

    Hook *hook.Hook
}

type Options struct {
    Writer  *writer.Writer
    Hook    *hook.Hook
    Trigger *Trigger
    IORetry Retry
}

func newJob(traceID string, reader reader.Reader, trigger *Trigger, writer *writer.Writer, retry Retry,
    hook *hook.Hook) Job {
    return Job{
        Log:               nil,
        TraceID:           traceID,
        TriggerConditions: trigger,
        Reader:            reader,
        Writer:            writer,
        IORetry:           retry,
        Hook:              hook,
    }
}

func DefaultRetry() Retry {
    retryInfo := lib.RetryInfo{
        Times:    3,
        Interval: 1 * time.Second,
    }
    return Retry{
        Trigger:            retryInfo,
        Job:                retryInfo,
        GetTransformStatus: retryInfo,
    }
}

func NewJob(traceID string, reader reader.Reader, opts ...Options) Job {

    var (
        write   *writer.Writer
        hooks   *hook.Hook
        trigger *Trigger
        retry   = DefaultRetry()
    )

    for _, opt := range opts {
        write = opt.Writer
        hooks = opt.Hook
        trigger = opt.Trigger
        retry = opt.IORetry
    }

    return newJob(traceID, reader, trigger, write, retry, hooks)
}

func (j *Job) InitJob() error {
    j.initLog()
    return nil
}

func (j Job) Trigger() (bool, error) {
    res, err := lib.Retry(j.IORetry.Trigger, func() (interface{}, error) {
        return j.TriggerConditions.Trigger()
    })
    return res.(bool), err
}

func (j Job) Run() error {
    return j.run()
}

func (j Job) runWithRetry() error {
    return lib.RetryOnlyReturnErr(j.IORetry.Job, func() error {
        return j.run()
    })
}

func (j *Job) run() error {
    defer func() {
        err := hook.ChainDefer(j.Hook.Defer)
        if err != nil {
            j.Log.Errorf("ChainDefer %v", err)
        }
    }()

    src, dst, err := hook.Chain(j.ReaderDst, j.ReaderDst, j.Hook.BeforeReader)
    if err != nil {
        return errors.Wrap(err, "j.Hook.BeforeReader")
    }

    err = j.Reader.Copy(src, dst)
    if err != nil {
        return errors.Wrap(err, "IReader.IO")
    }

    afterReader, err := hook.ChainDst(j.ReaderDst, j.Hook.AfterReader)

    if j.Writer != nil {
        err = j.Writer.Copy(afterReader, j.Writer.WriterDst)
        if err != nil {
            return errors.Wrap(err, "IWriter.Write")
        }
    }

    return nil
}
