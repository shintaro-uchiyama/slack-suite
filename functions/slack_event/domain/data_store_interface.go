package domain

type DataStoreInterface interface {
	Create(task Task) error
	Delete(timeStamp string) error
	Get(timeStamp string) (Task, error)
}
