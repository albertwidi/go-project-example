DROP TABLE IF EXISTS users;
CREATE TABLE users(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    hash_id varchar(8) NOT NULL,
    user_type smallint NOT NULL,
    user_status smallint NOT NULL,
    phone_number varchar(20) NOT NULL,
    email varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL,
    UNIQUE(phone_number)
);

DROP INDEX IF EXISTS idx_users_hash;
CREATE INDEX idx_users_hash ON users(hash_id);

DROP TABLE IF EXISTS users_bio;
CREATE TABLE users_bio(
    user_id bigint PRIMARY KEY,
    full_name varchar(60) NOT NULL,
    occupation varchar(30) NOT NULL,
    gender smallint NOT NULL,
    birthday date NOT NULL,
    avatar TEXT,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    updated_by varchar(36),
    is_test boolean NOT NULL
);

DROP TABLE IF EXISTS user_secrets;
CREATE TABLE user_secrets(
    id uuid PRIMARY KEY,
    user_id bigint NOT NULL,
    secret_key varchar(30) NOT NULL, -- secret key is unique per user_id
    secret_value varchar(100) NOT NULL, 
    created_at timestamp NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at timestamp,
    updated_by varchar(36),
    is_test boolean NOT NULL,
    UNIQUE(user_id, secret_key) 
);

DROP TABLE IF EXISTS registrations;
CREATE TABLE registrations(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid,
    user_type int NOT NULL,
    user_status int NOT NULL,
    ktp_id bigint NOT NULL,
    full_name varchar(60),
    birthdate date NOT NULL,
    email varchar(255) NOT NULL,
    phone_number varchar(20) NOT NULL,
    gender smallint NOT NULL,
    channel smallint, 
    device smallint,
    latitude varchar(20),
    longitude varchar(20),
    device_token varchar(200),
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL,
);