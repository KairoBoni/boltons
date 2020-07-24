package database

const (
	insertNFETotal = `INSERT INTO nfe_total(access_key, total)
	VALUES ($1, $2)`

	getNFETotal = `SELECT total FROM nfe_total
	WHERE access_key=$1`
)
