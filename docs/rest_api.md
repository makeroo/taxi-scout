# Taxi Scout REST API

Base url: `/api/v1/`

## Types

### EMAIL

Representation: STRING

### UUID

Representation: STRING

### DATETIME

Representation: STRING

Format: ISO 8601

### SCOUT_GROUP

Representation: OBJECT
* id: INT
* name: STRING

### INVITATION

Representation: OBJECT
* token: UUID
* email: EMAIL
* expires: DATETIME
* scout_group: SCOUT_GROUP

### ACCOUNT

Representation: OBJECT
* id: INT
* name: STRING
* email: EMAIL
* address: STRING


## Methods

### Accounts

**URL**: `/accounts`

**Method**: GET

Query accounts of a Scout Group.

Required permission: `member`

**Query parameters**:

* group INT

  A scout group.

*Response*: []ACCOUNT

**Errors**:

* 400

  Group parameter not found or not an int.

* 401

  Not authorized/authenticated.

* 403

  Forbidden: required permission is missing.


**Method**: POST

Create an account from an invitation and grant member permission.

Rquired permission: public

**Request**:

* invitation: UUID

  An invitation token.

**Response**:

* id: INT

  The account resulted from accepting the invitation.

* new_account: BOOL

  Whether a new account has been created or not.

* authenticated: BOOL

  Whether the invitation token has been used or not.

When the invitation token is used, a cookie named _ts_u is used. This cookie is used as authentication token
for next calls.

When cookie _ts_u is submitted, even if the invitation token is not found or expired, the call succeded
returning the account id of the logged user and new_account and authenticated both set to false.

**Errors**:

* 400

  Can't decode payload, missing invitation key or not an UUID

* 404 Invitation not found
  ```
  {
     error: 'not_found'
  }
  ```

* 410 Invitation expired
  ```
  {
     error: 'expired'
  }
  ```

**URL**: /account/:id

**Method**: GET

Query a specific account.

*Request*:
* id: the string "me" or an account ID

*Response*: ACCOUNT

Permission: a user can query its own account record.

**Errors**:

* 400

  Malformed :id parameter.

* 401

  Not authenticated.

* 403

  Forbidden.


**URL**: /account/:id/groups

**Method**: GET

*Response*: []SCOUT_GROUP

**URL**: /account/:id/group/:id/scouts

**Method**: GET

*Response*: []SCOUT




TODO FROM HERE



POST
request:
{
invitation: STRING (uuid token)
pwd: STRING
address: STRING
}

response:
{
id: INT
}

or 404
{
error: invitation_not_found
}

or 410
{
error: invitation_expired /* HTTP: gone */
}

// TODO: weak password?


/account/authenticate
POST:
{
login: STRING
pwd: STRING
}

Response: same as /account/:id


/api/v1/account/:id/scouts
GET

Response:
[
 {
  id: INT
  name: STRING
 }
]


Note: there is no specified order.

/api/v1/excursion/latest
GET

Response:
{
  detail: {
    id: INT,
    date: DATE (yyyy-mm-dd)
    from: TIME (hh:mm)
    to: TIME
    location: STRING
  },
  scouts: [
    {
      id: INT,
      name: STRING, // Note: not normalized
      partecipate: BOOL,
    }
  ],
  tutors: [
    {
      id: INT,
      name: STRING, // Note: not normalized
      rides: ENUM (0, 1, 2) 0 neither Out nor Return, 1 either Out or Return, 2 both Out and Return if needed
    }
  ],
  out: COORDINATION,
  return: COORDINATION,
}

COORDINATION:
{
  tutors: [
    id: INT,
    role: ENUM ('F', 'R')
    free_seats: INT,
    scouts: [INT]
  ],
  meetings: [
    {
      id: INT,
      taxi: INT,
      riders: [INT],
      point: STRING,
      time: TIME
    }
  ]
}





### Invitations

**URL**: `/invitations`

**Method**: POST

Create a new invitation.

Required role: `excursion_manager`

*Request*:

* email: EMAIL

  Email address of invitation receiver.

* scout_group: INT

  Scout group to be joined.

*Response*:

* token: UUID

  New inviation identifier.

* expires: DATETIME
