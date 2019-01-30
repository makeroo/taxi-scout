package storage

var SqlQueries = map[string]string{
	"check_permission": "SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?",
	"grant": "INSERT INTO account_grant (permission_id, account_id, group_id) VALUES (?, ?, ?)",

	"fetch_invitation": `
   SELECT i.email, i.created_on, i.group_id,
          a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?
`,

	"delete_invitation": "DELETE FROM invitation WHERE token = ?",

	"create_account_from_invitation": "INSERT INTO account (email) SELECT email FROM invitation WHERE token = ?",

	"query_accounts": `
SELECT a.id, a.name, a.email, a.address
  FROM account a
  JOIN account_grant g ON g.account_id = a.id
 WHERE g.group_id = ?`,

	"query_account": "SELECT id, name, email FROM account WHERE id = ?",
	//"insert_account": "INSERT INTO account (name, email, password, address) VALUES ( ?, ?, ?, ? )",
	"account_credentials": "SELECT id, password FROM account WHERE email = ?",
}
