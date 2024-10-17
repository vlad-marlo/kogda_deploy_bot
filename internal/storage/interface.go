package storage

type ChatRepository interface {
	GetConsumersIDs() ([]int64, error)
}

type Storage interface {
	Chat() ChatRepository
}
