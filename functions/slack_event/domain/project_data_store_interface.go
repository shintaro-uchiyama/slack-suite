package domain

type ProjectDataStoreInterface interface {
	GetByChannel(channel string) (Project, error)
}
