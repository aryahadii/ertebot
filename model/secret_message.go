package model

type SecretMessage struct {
	Message          string
	SenderID         string
	SenderUsername   string
	ReceiverID       string
	ReceiverUsername string
	SendEpoch        int64
	SeenEpoch        int64
}

type Message struct {
	Message          string
	SenderID         string
	SenderUsername   string
	ReceiverID       string
	ReceiverUsername string
	SendEpoch        int64
	SeenEpoch        int64
}

type SecretMessageNewFirst []SecretMessage

func (m SecretMessageNewFirst) Len() int {
	return len(m)
}

func (m SecretMessageNewFirst) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m SecretMessageNewFirst) Less(i, j int) bool {
	return m[i].SendEpoch > m[j].SendEpoch
}

type ThreadNewFirst [][]SecretMessage

func (t ThreadNewFirst) Len() int {
	return len(t)
}

func (t ThreadNewFirst) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ThreadNewFirst) Less(i, j int) bool {
	return t[i][0].SendEpoch > t[j][0].SendEpoch
}
