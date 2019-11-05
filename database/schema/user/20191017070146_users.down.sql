DROP TABLE IF EXISTS users;
CREATE TABLE users(
    id uuid PRIMARY KEY,
    hash_id varchar(6) NOT NULL,
    user_type smallint NOT NULL,
    user_status smallint NOT NULL,
    phone_number varchar(20) NOT NULL,
    email varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL,
    UNIQUE(phone_number)
);