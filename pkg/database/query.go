package database

const (
	insertNFEAmount = `INSERT INTO nfe_amount(access_key, amount)
	VALUES ($1, $2);`

	getNFEAmount = `SELECT amount FROM nfe_amount
	WHERE access_key=$1;`
)
