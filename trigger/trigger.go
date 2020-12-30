package trigger

type Trigger interface {
	CanStart() (bool, error)
}
