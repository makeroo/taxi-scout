package storage

import "github.com/kataras/iris/core/errors"

type Datastore interface {
	QueryInvitationToken (token string) (*Account, bool, error)

	QueryAccounts(group int32, userId int32) ([]*Account, error)

	QueryAccount(id int32) (*Account, error)

	AuthenticateAccount(email string, pwd string) (int32, error)

	UpdateAccountPassword(id int32, oldPwd string, newPwd string) error

	AccountGroups(accountId int32) ([]*ScoutGroup, error)

	AccountScouts(accountId int32) ([]*Scout, error)

	AccountUpdate(account *Account) error

	InsertOrUpdateScout(scout Scout, tutorId int32) (int32, error)
}

var IdOverflow = errors.New("id_overflow")

var UnknownQuery = errors.New("unknown_query")
