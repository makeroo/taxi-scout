package storage

import "errors"

// NoRequestingUser is an invalid key to indicate no user.
const NoRequestingUser = int32(-1)

// Datastore is TaxiScount storage abstraction.
type Datastore interface {
	CreateInvitationForExistingMember(email string) (*Invitation, error)

	// use NoRequestingUser if no user is authenticated, otherwise his/her id
	QueryInvitationToken(token string, requestingUser int32) (invitedAccount *Account, accountAlreadyExisted bool, scoutGroupID int32, joinedGroup bool, err error)

	QueryAccounts(group int32, userID int32) ([]*Account, error)

	QueryAccount(id int32) (*Account, error)

	AuthenticateAccount(email string, pwd string) (int32, error)

	UpdateAccountPassword(id int32, newPwd string) error

	AccountGroups(accountID int32) ([]*ScoutGroup, error)

	AccountScouts(accountID int32) ([]*Scout, error)

	AccountUpdate(account *Account) error

	InsertOrUpdateScout(scout Scout, tutorID int32) (int32, error)

	RemoveScout(scoutID int32, tutorID int32) error
}

// ErrIDOverflow is returned when account creation fails due to overflow of ID field.
var ErrIDOverflow = errors.New("id_overflow")

// ErrUnknownQuery denotes a bug in SQLDatastore.
var ErrUnknownQuery = errors.New("unknown_query")
