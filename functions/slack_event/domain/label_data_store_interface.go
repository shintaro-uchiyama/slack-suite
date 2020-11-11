package domain

type LabelDataStoreInterface interface {
	GetByReaction(channel string, reaction string) (Label, error)
}
