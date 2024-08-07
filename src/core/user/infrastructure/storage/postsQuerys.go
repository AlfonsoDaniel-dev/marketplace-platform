package Userstorage

const sqlCreatePost = `INSERT INTO posts(id, creator_id, content_path, title, content, created_at) VALUES ($1, $2, $3, $4, $5)`

const sqlUpdatePostTitle = `UPDATE posts SET title =$1, updated_at= $2 WHERE id = $3`

const sqlUpdatePostContent = `UPDATE posts SET content = $1, updated_at = $2 WHERE creator_id = $3`
