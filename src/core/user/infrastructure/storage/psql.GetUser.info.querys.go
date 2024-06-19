package Userstorage

const sqlGetUserNameByEmail = `SELECT user_name FROM users WHERE email = $1`

const sqlInsertUserAddress = `INSERT INTO user_addresses(id, user_id, street, city, state, postal_code, country, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

const sqlAddAddressOnUseTable = `INSERT INTO users (address_id) SELECT user_id FROM user_addresses WHERE address_id = $1`

const sqlCheckUserTsvIsTrue = `SELECT two_steps_verification FROM users WHERE email = $1`
