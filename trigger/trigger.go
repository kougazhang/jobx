package trigger

type Trigger interface {
    CanStart(p ...interface{}) (bool, error)
}
