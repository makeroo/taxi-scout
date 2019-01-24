package storage

var SqlQueries = map[string]string{
	"fetch_invitation": `
   SELECT i.email, i.created_on,
          g.id, g.name,
          a.id, a.name, a.address, a.verified_email
     FROM invitation i
     JOIN scout_group g ON i.group_id = g.id
LEFT JOIN account a ON i.account_id = a.id
    WHERE i.token = ?
`,

	"query_accounts": "SELECT id, name, email FROM account",
	"query_account": "SELECT id, name, email FROM account WHERE id = ?",
	"insert_account": "INSERT INTO account (name, email, password, address) VALUES ( ?, ?, ?, ? )",
	"account_credentials": "SELECT id, password FROM account WHERE email = ?",
}
