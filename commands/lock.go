package commands

type Lock struct {
	Unlock bool `json:"unlock,omitempty"`
}

func _(lock bool) *Lock {
	return &Lock{
		Unlock: lock,
	}
}
