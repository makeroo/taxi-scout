package storage

type Account struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	VerifiedEmail bool `json:"verified_email"`
	Address string `json:"address"`
}

type AccountWithCredentials struct {
	*Account
	Pwd string `json:"pwd"`
}

func NewAccount() *Account {
	return &Account{0, "", "", false, ""}
}

func NewAccountWithCredentials() *AccountWithCredentials {
	return &AccountWithCredentials{NewAccount(), ""}
}
