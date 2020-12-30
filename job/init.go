package job

import log "github.com/sirupsen/logrus"



func (j *Job) initLog() {
	j.Log = log.WithField("TraceID", j.TraceID)
}
