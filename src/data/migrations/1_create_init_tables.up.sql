
CREATE TABLE ac_keys (
    id VARCHAR(255) PRIMARY KEY,
    premium BOOLEAN DEFAULT FALSE, -- Use BOOLEAN instead of bool
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NULL -- Specify WITH TIME ZONE for consistency
);

CREATE TABLE active_devices (
    id SERIAL PRIMARY KEY,
    device_name VARCHAR(255) NULL,
    session_id VARCHAR(255) NOT NULL, -- Use VARCHAR instead of string
    ip VARCHAR(255) NULL,
    ac_keys_id VARCHAR(255) NOT NULL, -- Define the ac_keys_id column explicitly
    FOREIGN KEY (ac_keys_id) REFERENCES ac_keys(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (session_id, ac_keys_id) -- Add this line to enforce the composite unique constraint
);
