package migrations

const sqlMigrateUuidExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`

const sqlMigrateUserTable = `CREATE TABLE IF NOT EXISTS users(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    first_name VARCHAR(70) NOT NULL,
    last_name VARCHAR(70) NOT NULL,
    user_name VARCHAR(70) NOT NULL,
    profile_picture_id UUID NOT NULL,
    media_repository_path VARCHAR(250),
    collections_table_id uuid,
    biography VARCHAR(140) NOT NULL,
    age INT NOT NULL,
    email VARCHAR(60) NOT NULL,
    password VARCHAR(350) NOT NULL,
    two_steps_verification BOOLEAN NOT NULL DEFAULT false,
    address_id uuid,
    
    
    created_at INT,
    updated_at INT,
    deleted_at INT,
    
    CONSTRAINT users_id_pk PRIMARY KEY (id),
    CONSTRAINT users_age_ck CHECK (age >= 18),
    CONSTRAINT users_first_and_last_name_uq UNIQUE (first_name, last_name),
    CONSTRAINT users_email_uq UNIQUE (email),
    CONSTRAINT users_user_name UNIQUE (user_name)
)`

const sqlMigrateUserCollectionTable = `CREATE TABLE IF NOT EXISTS user_collections(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
)`

const sqlUsersRelationWithUserCollectionsTable = `ALTER TABLE users ADD CONSTRAINT users_collections_table_id_fk FOREIGN KEY (collections_table_id) REFERENCES user_collections (id) ON UPDATE CASCADE`

const sqlUserCollectionsRelationWithUsersTable = `ALTER TABLE user_collections ADD CONSTRAINT user_collections_user_id FOREIGN KEY (user_id) REFERENCES users_id ON UPDATE CASCADE ON DELETE RESTRICT`

const sqlMigrateCollectionsTable = `CREATE TABLE IF NOT EXISTS collections(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    user_collections_id uuid NOT NULL,
    collection_name VARCHAR(100) NOT NULL,
    collection_description VARCHAR(140),
    
    collection_path VARCHAR(250) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    
    CONSTRAINT collections_id_pk PRIMARY KEY (id),
    CONSTRAINT collections_name_uq UNIQUE (collection_name),
    CONSTRAINT collections_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
    CONSTRAINT collections_user_collections_id FOREIGN KEY (user_collections_id) REFERENCES user_collections (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT collections_collection_path_ck CHECK (collection == '')
    
)`

const sqlMigrateImagesTable = `CREATE TABLE IF NOT EXISTS images(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    collection_id uuid NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_path VARCHAR(150) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
        
    CONSTRAINT images_id_pk PRIMARY KEY (id),
    CONSTRAINT images_file_name_ck CHECK (file_name = ''),
    CONSTRAINT images_file_path_ck CHECK (file_path = ''),
    CONSTRAINT images_user_id_fk FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
)`

const sqlMigrateProfilePicturesId = `CREATE TABLE IF NOT EXISTS profilePictures(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    user_Repository_path VARCHAR(140) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_path VARCHAR(150) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    
    CONSTRAINT profilePictures_id_pk PRIMARY KEY (id),
    CONSTRAINT profilePictures_file_name_ck CHECK (file_name = ''),
    CONSTRAINT profilePictures_file_path_ck CHECK (file_path = ''),
    CONSTRAINT profilePictures_user_id_fk FOREIGN KEY (user_id)
    REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT
)`

const sqlMigrateUserAddressesTable = `CREATE TABLE IF NOT EXISTS user_addresses(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    
    street VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100),
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),

	CONSTRAINT user_addresses_id_pk PRIMARY KEY (id),
	CONSTRAINT user_addresses_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE,
	CONSTRAINT user_addresses_user_id_uq UNIQUE (user_id)
)`

// makes the foreign key with user_addresses table on users table
const sqlAddConstraintForUserWithAddress = `
	ALTER TABLE users 
    ADD CONSTRAINT users_address_id_fk
    FOREIGN KEY (address_id)
    REFERENCES user_addresses (id)
	ON UPDATE CASCADE
`

/*

const sqlMigrateUserRatingTable = `CREATE TABLE IF NOT EXISTS user_ratings (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_rated_id UUID NOT NULL,
    who_is_rating_id UUID NOT NULL,
    rating INT,

    CONSTRAINT user_ratings_pk PRIMARY KEY (id),
    CONSTRAINT user_ratings_rating_ck CHECK (rating >=0 AND rating <= 5),
    CONSTRAINT user_ratings_user_rated_id_fk FOREIGN KEY (user_rated_id)
    REFERENCES users (id)
    ON UPDATE CASCADE,
    CONSTRAINT user_ratings_who_is_raiting_id_fk FOREIGN KEY (who_is_rating_id)
    REFERENCES users (id)
    ON UPDATE CASCADE,
)`
/*
*/

/*

const sqlMigrateUserCommentReviewTable = `CREATE TABLE IF NOT EXISTS user_comment_reviews (
    comment_review_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    -- reviewer_id is who review another user
    reviewer_id UUID NOT NULL,
    -- user_reviewed_id is who it's been reviewed
    user_reviewed_id UUID NOT NULL,

    title VARCHAR(80) NOT NULL,
    message VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT user_comment_reviews_id_pk PRIMARY KEY (id),
    CONSTRAINT user_comment_reviews_reviewer_id_fk FOREIGN KEY (reviewer_id)
	REFERENCES users (id)
    ON UPDATE CASCADE,
    CONSTRAINT user_comment_reviews_user_reviewed_id_fk FOREIGN KEY (user_reviewed_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
)` */

/*
const sqlMigrateShopThemesTable = `CREATE TABLE products_and_shop_themes(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    theme_name VARCHAR(140) NOT NULL,
    description VARCHAR(140) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,

    deleted_at TIMESTAMP,

    CONSTRAINT products_and_shop_themes_id_pk PRIMARY KEY (id),
    CONSTRAINT products_and_shop_themes_theme_name UNIQUE (theme_name),
    CONSTRAINT products_and_shop_themes_description UNIQUE (description)
)` */
/*
const sqlMigrateSellersAccounts = `CREATE TABLE IF NOT EXISTS sellers_accounts(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
)` */
/*
const sqlMigrateShopsTable = `CREATE TABLE IF NOT EXISTS shops(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    seller_id UUID NOT NULL,
    shop_name VARCHAR(80) NOT NULL,
    shop_email VARCHAR(80) NOT NULL,

    shop_theme_name VARCHAR(140) NOT NULL,

)` */

const sqlMigrateProductCategoriesTable = `CREATE TABLE IF NOT EXISTS product_categories(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    theme_name VARCHAR(100),
    description VARCHAR(150),
    
    CONSTRAINT product_categories_id_pk PRIMARY KEY (id)
)`

/*
const sqlMigrateShopsProductsTable = `CREATE TABLE IF NOT EXISTS shop_products(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    seller_id uuid NOT NULL,

    name VARCHAR(60) NOT NULL,
    category_id UUID NOT NULL,
    description VARCHAR(140),
    stock INT NOT NULL,
    price REAL NOT NULL DEFAULT 0.0,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT products_id_pk PRIMARY KEY (id),

    CONSTRAINT products_shop_id_fk
    FOREIGN KEY (seller_id)
    REFERENCES shops (id)
    ON UPDATE CASCADE,

    CONSTRAINT products_category_id_fk
    FOREIGN KEY (category_id)
    REFERENCES product_categories(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
)` */

const sqlMigrateUserProductsTable = `CREATE TABLE IF NOT EXISTS products(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    seller_id uuid NOT NULL,
    
    name VARCHAR(60) NOT NULL,
    category_id UUID NOT NULL,
    description VARCHAR(140),
    stock INT NOT NULL,
    price REAL NOT NULL DEFAULT 0.0,
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT products_id_pk PRIMARY KEY (id),
    
    CONSTRAINT products_seller_id_fk 
    FOREIGN KEY (seller_id) 
    REFERENCES sellers (id)
    ON UPDATE CASCADE,
    
    CONSTRAINT products_category_id_fk
    FOREIGN KEY (category_id)
    REFERENCES product_categories(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
)`

const sqlMigratePurchaseProductTable = `CREATE TABLE IF NOT EXISTS purchase_products(
    order_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT '',
    product_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    total REAL NOT NULL DEFAULT 0.0
    
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    cancel_at TIMESTAMP,
)`

/*
const sqlMigrateOrderTable = `CREATE TABLE IF NOT EXISTS orders(
    order_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    seller_id uuid NOT NULL,
    shop_id_uuid uuid,
    status VARCHAR(10) NOT NULL DEFAULT 'IN_PROGRESS',

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)`

const sqlMigrateOrderItemsTable = `CREATE TABLE IF NOT EXISTS order_items(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    product_id uuid NOT NULL,
    quantity INT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    delted_at TIMESTAMP,

    CONSTRAINT order_items_id_pk PRIMARY KEY (id),

    CONSTRAINT order_items_order_id_fk FOREIGN KEY (order_id)
    REFERENCES orders (id)
    ON DELETE CASCADE,

    CONSTRAINT order_items_product_id_fk FOREIGN KEY (product_id)
    REFERENCES products (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
)`
*/
