package storage

import "github.com/kataras/iris/core/errors"

type Datastore interface {
	CheckPermission (userId int32, groupId int32, permId int32) error;

	QueryInvitationToken (token string) (*Account, bool, error)

	QueryAccounts(group int32) ([]*Account, error)

	QueryAccount(id int32) (*Account, error)

//	InsertAccount(*AccountWithCredentials) (int32, error)

	AuthenticateAccount(email string, pwd string) (int32, error)

	UpdateAccountPassword(id int32, oldPwd string, newPwd string) error
}

var IdOverflow = errors.New("id_overflow")

var UnknownQuery = errors.New("unknown_query")
