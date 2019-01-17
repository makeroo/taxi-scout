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
* verified_email: BOOL
* address: STRING


## Methods

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


**URL**: `/invitation/:token`

**Method**: GET

Retrieve invitation details.

Public method.

*Response*:

* type: ENUM 'account', 'invitation'

* invitation: INVITATION

  If type is 'invitation'

* account: ACCOUNT

  If type is 'account'


Errors:

* 404

  ```
  {
     error: 'not_found'
  }
  ```

 * 410

   ```
   {
      error: 'expired'
   }
   ```


/accounts
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

/account/:id
GET

Response:
{
id: INT
name: STRING
email: STRING
verified_email: BOOL
address: STRING
}

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
