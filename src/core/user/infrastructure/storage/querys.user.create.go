package Userstorage

// sqLCreateUser inserts user data for register
const sqlCreateUser = `INSERT INTO users (id, first_name, last_name, user_name, biography, age, email, password, two_steps_verification, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

const sqlGetUserIdByEmail = `SELECT id FROM users WHERE email = $1`

// inserts new address on user_addresses table
const sqlInsertAddressOnAddressTable = `INSERT INTO user_addresses (id, user_id, street, city, state, postal_code, country, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

const sqlChangeUserName = `UPDATE users SET user_name = $1, updated_at = $2 WHERE email = $3`

const sqlChangeUserPassword = `UPDATE users SET password = $1, updated_at = $2 WHERE password = $3 AND email = $4`

const sqlChangeUserFirstName = `UPDATE users SET first_name = $1, updated_at = $2 WHERE email = $3`

const sqlChangeUserLastName = `UPDATE users SET last_name = $1, updated_at = $2 WHERE email = $3`

const sqlChangeUserEmail = `UPDATE users SET email = $1, updated_at = $2 WHERE id = $3`

const sqlChangeTsvStatus = `UPDATE users SET two_steps_verification = $1 WHERE email = $2`
