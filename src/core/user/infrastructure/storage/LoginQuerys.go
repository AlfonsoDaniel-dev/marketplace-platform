package Userstorage

const sqlGetHashedPasswordFromEmail = `SELECT password FROM users WHERE email = $1`

const sqlLoginVerifyEmailExists = `SELECT email FROM users WHERE email = $1`
