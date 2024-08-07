package migrations

const sqlMigratePostsTable = `CREATE TABLE IF NOT EXISTS posts(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    creator_id uuid NOT NULL,
    content_path VARCHAR (250) NOT NULL,
    
    title VARCHAR(100) NOT NULL,
    content VARCHAR(340) NOT NULL,
    
    created_at int NOT NULL,
    updated_at int,
    
    CONSTRAINT posts_id_pk PRIMARY KEY (id),
    CONSTRAINT posts_creator_id_fk FOREIGN KEY (creator_id)
    REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE
)`
