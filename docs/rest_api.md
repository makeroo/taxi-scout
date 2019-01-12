/api/v1/

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
