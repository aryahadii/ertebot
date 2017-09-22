package model

type Person struct {
	UserID              string
	FirstName           string
	LastName            string
	Username            string
	HashID              string
	LastUseEpoch        int64
	LastInboxCheckEpoch int64
}

const (
	HashIDLength = 10
)
