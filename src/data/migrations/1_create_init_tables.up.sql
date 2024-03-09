CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


SET TIMEZONE="Iran";

-- CREATE TYPE user_type_enum AS ENUM ('paid', 'vip', 'admin');

CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    email VARCHAR (255) NOT NULL UNIQUE,
    user_password VARCHAR (255) NOT NULL,
    user_type user_type_enum NOT NULL DEFAULT 'paid',
    user_attrs JSONB NULL
);

CREATE TABLE active_devices (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
