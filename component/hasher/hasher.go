package hasher

type Hasher interface {
	Hash(data string) string
	CheckPassword(hashedPassword, password string) error
}
