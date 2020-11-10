package domain

type ZubeInterface interface {
	Create(task Task) (int, error)
	Delete(cardID int) error
}
