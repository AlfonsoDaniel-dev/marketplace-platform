package Userstorage

// sqLCreateUser inserts user data for register
const sqlCreateUser = `INSERT INTO users (id, first_name, last_name, user_name, biography, age, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

const sqlGetUserIdByEmail = `SELECT id FROM users WHERE email = $1`

// inserts new address on user_addresses table
const sqlInsertAddressOnAddressTable = `INSERT INTO user_addresses (id, user_id, street, city, state, postal_code, country, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
