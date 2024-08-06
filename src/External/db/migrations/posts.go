package migrations

const MigratePostsTable = `CREATE TABLE IF NOT EXISTS posts (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    creator_id uuid NOT NULL,
    user_posts_path VARCHAR(255) NOT NULL,
    content_path VARCHAR (250) NOT NULL,
    
    title VARCHAR(100) NOT NULL,
    content VARCHAR(340) NOT NULL,
    
    created_at int NOT NULL,
    updated_at int NOT NULL,
    
    CONSTRAINT posts_id_pk PRIMARY KEY (id),
    CONSTRAINT posts_creator_id_fk FOREIGN KEY (creator_id)
    REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT user_posts_path_fk FOREIGN KEY (user_posts_path)
    REFERENCES users (posts_directory) ON DELETE RESTRICT ON UPDATE CASCADE
)`
