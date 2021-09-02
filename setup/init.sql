
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50),
    phone_number VARCHAR(50) UNIQUE NOT NULL,
    is_verified BOOL DEFAULT FALSE
);

CREATE TABLE otp (
    value VARCHAR(50),
    phone_number VARCHAR(50) UNIQUE NOT NULL,
    expiry timestamp
)
