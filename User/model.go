package User

type User struct {
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (u User) UniqKey() interface{} {
	return u.Name
}