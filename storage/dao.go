package storage

type Datastore interface {
	QueryInvitationToken (token string) (*Invitation, *Account, error)

	QueryAccounts() ([]*Account, error)

	QueryAccount(id int32) (*Account, error)

	InsertAccount(*AccountWithCredentials) (int32, error)

	AuthenticateAccount(email string, pwd string) (int32, error)

	UpdateAccountPassword(id int32, oldPwd string, newPwd string) error
}

const IdOverflow = "id_overflow"

const UnknownQuery = "unknown_query"