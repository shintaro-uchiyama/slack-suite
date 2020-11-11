package domain

type DataStoreInterface interface {
	Create(task Task) error
	Delete(channel string, timeStamp string) error
	Get(channel string, timeStamp string) (Task, error)
}
