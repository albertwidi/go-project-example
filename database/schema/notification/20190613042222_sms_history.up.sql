-- extension for uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS notification_types;
CREATE TABLE notification_types (
    id int PRIMARY KEY,
    notification_type_name varchar(30) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL
);

-- initial data of notification type
INSERT INTO notification_types(id, notification_type_name, created_at, is_test) VALUES(1, 'sms', CURRENT_TIMESTAMP, false);
INSERT INTO notification_types(id, notification_type_name, created_at, is_test) VALUES(2, 'email', CURRENT_TIMESTAMP, false);
INSERT INTO notification_types(id, notification_type_name, created_at, is_test) VALUES(3, 'push-message', CURRENT_TIMESTAMP, false);

DROP TABLE IF EXISTS notification_providers;
CREATE TABLE notification_providers (
    id int PRIMARY KEY,
    notification_provider_type int NOT NULL,
    notification_provider_name VARCHAR(30) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL
);

-- initial data for notifications provider
INSERT INTO notification_providers(id, notification_provider_type, notification_provider_name, created_at, is_test) VALUES(1, 1, 'nexmo', CURRENT_TIMESTAMP, false);

DROP TABLE IF EXISTS notification_purpose;
CREATE TABLE notification_purpose (
    id int PRIMARY KEY,
    notification_purpose_name varchar(30) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_test boolean NOT NULL
);

-- initial data for notification purposes
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(1, 'system-update', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(2, 'authentication-otp', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(3, 'authentication-payment', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(4, 'authentication-withdraw', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(5, 'promotion', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(6, 'reminder', CURRENT_TIMESTAMP, false);
INSERT INTO notification_purpose(id, notification_purpose_name, created_at, is_test) VALUES(7, 'misc', CURRENT_TIMESTAMP, false);

DROP TABLE IF EXISTS sms_history;
CREATE TABLE sms_history (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    -- provider of sms. for example nexmo
    notification_provider int NOT NULL,
    notification_provider_name varchar(30) NOT NULL,
    -- message id from provider
    message_id varchar(50) NOT NULL,
    message_count int NOT NULL,
    to_phone_number VARCHAR(20) NOT NULL,
    from_phone_number VARCHAR(20) NOT NULL,
    -- internal delivery status
    delivery_status int NOT NULL,
    -- delivery status from provider, may vary
    provider_delivery_status VARCHAR(20) NOT NULL,
    -- sms type to differentiate sms purposes, for example opt, marketing, etc
    message_type int NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    is_test boolean NOT NULL
);

DROP TABLE IF EXISTS email_history;
CREATE TABLE email_history (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    -- provider of email
    notification_provider int NOT NULL,
    notification_provider_name varchar(30) NOT NULL,
    -- email specific
    email_title varchar(200) NOT NULL,
    email_body text NOT NULL,
    -- status
    provider_delivery_status int NOT NULL,
    delivery_status int NOT NULL,
    -- metadata
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_text boolean NOT NULL
);

DROP TABLE IF EXISTS user_notifications;
CREATE TABLE user_notifications(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id bigint NOT NULL,
    -- provider of notification, sms, email, push
    notification_provider_type int NOT NULL,
    -- provider of notification, nexmo, sendgrid, etc
    notification_provider_id int NOT NULL,
    -- send id that provided by the notification provider
    notification_provider_send_id varchar(30) NOT NULL,
    -- purpose of the notification. marketing, information, etc
    notification_purpose int NOT NULL, 
    -- if notification is webpage, use this
    notification_is_webpage boolean,
    -- status of notification
    notification_status int NOT NULL,
    -- minimum information about the notification
    notification_title varchar(100) NOT NULL,
    notification_message varchar(255) NOT NULL,
    -- a state to set whether notification will be showed or not
    notification_show boolean NOT NULL,
    -- to give a hint whether a notification has detail or not
    notification_has_detail boolean NOT NULL,
    -- for read or unread flag
    notification_read boolean NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    is_deleted boolean NOT NULL,
    is_test boolean NOT NULL
);

DROP INDEX IF EXISTS idx_user_notifications;
CREATE INDEX idx_user_notifications ON user_notifications(user_id, created_at);

DROP TABLE IF EXISTS user_notifications_detail;
CREATE TABLE user_notifications_detail(
    notification_id uuid PRIMARY KEY,
    user_id bigint NOT NULL,
    notification_body text NOT NULL,
    notification_web_link text NOT NULL,
    created_at timestamp NOT NULL,
    is_deleted boolean NOT NULL,
    is_test boolean NOT NULL
);