package commands

type Lock struct {
	Lock bool `json:"lock,omitempty"`
}

func _(lock bool) *Lock {
	return &Lock{
		Lock: lock,
	}
}
