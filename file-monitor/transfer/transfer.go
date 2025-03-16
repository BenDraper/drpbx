package transfer

type Transfer interface {
	Create(entry string) error
	Update(entry string) error
	Delete(entry string) error
}
