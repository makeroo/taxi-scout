package storage

type Account struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountWithCredentials struct {
	*Account
	Pwd string `json:"pwd"`
}

func NewAccount() *Account {
	return &Account{0, "", ""}
}

func NewAccountWithCredentials() *AccountWithCredentials {
	return &AccountWithCredentials{NewAccount(), ""}
}
