package storage

// SQLQueries is a collection to all static SQL queries used by SQLDatastore.
var SQLQueries = map[string]string{
	"check_permission": "SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?",
	"grant":            "INSERT INTO account_grant (permission_id, account_id, group_id) VALUES (?, ?, ?)",

	"fetch_invitation": `
   SELECT i.email, i.created_on, i.group_id,
          a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`,

	"delete_invitation": "DELETE FROM invitation WHERE token = ?",
	"create_invitation_for_existing_member": `
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`,

	"create_account_from_invitation": "INSERT INTO account (email) SELECT email FROM invitation WHERE token = ?",

	"query_accounts": `
SELECT a.id, a.name, a.email, a.address
  FROM account a
  JOIN account_grant g ON g.account_id = a.id
 WHERE g.group_id = ?`,

	"query_account": "SELECT id, name, email, address FROM account WHERE id = ?",
	//"insert_account": "INSERT INTO account (name, email, password, address) VALUES ( ?, ?, ?, ? )",
	"update_account":      "UPDATE account SET name = ?, address = ? WHERE id = ?",
	"account_credentials": "SELECT id, password FROM account WHERE email = ?",

	"account_groups": `
SELECT g.id, g.name
  FROM scout_group g
  JOIN account_grant ag ON ag.group_id = g.id
 WHERE ag.account_id = ? AND ag.permission_id = ?
`,

	"account_scouts": `
SELECT s.id, s.name, s.group_id
  FROM scout s
  JOIN tutor_scout t ON t.scout_id = s.id
 WHERE t.tutor_id = ?
`,

	"add_scout":          "INSERT INTO tutor_scout (tutor_id, scout_id) VALUES (?, ?)",
	"remove_scout_tutor": "DELETE FROM tutor_scout WHERE scout_id = ? AND tutor_id = ?",
	"count_tutors":       "SELECT count(*) FROM tutor_scout WHERE scout_id = ?",

	"check_if_tutor": "SELECT count(*) FROM tutor_scout WHERE scout_id=? AND tutor_id=?",

	"insert_scout": "INSERT INTO scout (name, group_id) VALUES (?, ?)",
	"update_scout": "UPDATE scout SET name = ?, group_id = ? WHERE id = ?",
	"remove_scout": "DELETE FROM scout WHERE id = ?",
}
