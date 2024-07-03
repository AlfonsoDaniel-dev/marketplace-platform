package migrations

const sqlMigrateCollectionsTable = `CREATE TABLE IF NOT EXISTS collections(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    collection_name VARCHAR(100) NOT NULL,
    collection_description VARCHAR(140),
    
    collection_path VARCHAR(250) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    
    CONSTRAINT collections_id_pk PRIMARY KEY (id),
    CONSTRAINT collections_name_uq UNIQUE (collection_name),
    CONSTRAINT collections_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    CONSTRAINT collections_collection_path_ck CHECK (collection_path != '')
)`

const sqlMigrateImagesTable = `CREATE TABLE IF NOT EXISTS images(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    collection_id uuid,
	collection_path VARCHAR(250),
    user_repository_path VARCHAR(80) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_extension VARCHAR(8) NOT NULL,
    file_path VARCHAR(150) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
        
    CONSTRAINT images_id_pk PRIMARY KEY (id),
    CONSTRAINT images_collection_path_ck CHECK (collection_path != ''),
    CONSTRAINT images_file_name_ck CHECK (file_name != ''),
    CONSTRAINT images_file_path_ck CHECK (file_path != ''),
    CONSTRAINT images_user_id_fk FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
    CONSTRAINT images_collection_id_fk FOREIGN KEY (collection_id)
    REFERENCES collections (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
)`
