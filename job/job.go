package job

import (
    "github.com/kougazhang/jobx/io"
    "github.com/kougazhang/jobx/lib"
    "github.com/kougazhang/jobx/plugins"
    "github.com/kougazhang/jobx/reader"
    "github.com/kougazhang/jobx/writer"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "time"
)

type Job struct {
    Log     *log.Entry
    TraceID string

    // 触发任务
    TriggerConditions Trigger

    // 读取
    reader.Reader

    // 输出
    writer.Writer

    Retry Retry

    Plugins plugins.Plugins
}

func NewJob(traceID string, reader reader.Reader, trigger Trigger, writer writer.Writer, retry Retry,
    plugins plugins.Plugins) Job {
    return Job{
        Log:               nil,
        TraceID:           traceID,
        TriggerConditions: trigger,
        Reader:            reader,
        Writer:            writer,
        Retry:             retry,
        Plugins:           plugins,
    }
}

func NewJobWithDefault(traceID string, reader reader.Reader, trigger Trigger, writer writer.Writer,
    plugins plugins.Plugins) Job {

    retryInfo := lib.RetryInfo{
        Times:    3,
        Interval: 1 * time.Second,
    }
    retry := Retry{
        Trigger:            retryInfo,
        Job:                retryInfo,
        GetTransformStatus: retryInfo,
    }
    return NewJob(traceID, reader, trigger, writer, retry, plugins)
}

func (j *Job) InitJob() error {
    j.initLog()
    return nil
}

func (j Job) Trigger() (bool, error) {
    res, err := lib.Retry(j.Retry.Trigger, func() (interface{}, error) {
        return j.TriggerConditions.Trigger()
    })
    return res.(bool), err
}

func (j Job) Run() error {
    return j.run()
}

func (j Job) runWithRetry() error {
    return lib.RetryOnlyReturnErr(j.Retry.Job, func() error {
        return j.run()
    })
}

func (j *Job) run() error {
    err := j.Reader.Copy(j.ReaderSrc, j.ReaderDst)
    if err != nil {
        return errors.Wrap(err, "IReader.IO")
    }

    afterReaders, err := io.ChainDst(j.ReaderDst, j.Plugins.AfterReaders)

    if j.Writer.IWriter != nil {
        err = j.Writer.Copy(afterReaders, j.WriterDst)
        if err != nil {
            return errors.Wrap(err, "IWriter.Write")
        }
    }

    return nil
}
