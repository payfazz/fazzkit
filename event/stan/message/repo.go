package message

//Repository persistence for event message
type Repository interface {
	Save(message Message) error
	Load(ID []byte) (*Message, error)
	Close()
}
