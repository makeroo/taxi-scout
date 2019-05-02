package storage

// Account collects minimal informations about TaxiScout user.
type Account struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

/*
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
*/
