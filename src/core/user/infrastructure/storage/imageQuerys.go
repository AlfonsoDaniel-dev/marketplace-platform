package Userstorage

const sqlInsertImageData = `INSERT INTO images (id, user_id, user_repository_path, file_name, file_extension, file_path, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

const sqlInsertRepositoryPathOnUser = `UPDATE users SET media_repository_path = $1 WHERE id = $2`
