package model

type Person struct {
	UserID              string
	FirstName           string
	LastName            string
	Username            string
	LastUseEpoch        int64
	LastInboxCheckEpoch int64
}
