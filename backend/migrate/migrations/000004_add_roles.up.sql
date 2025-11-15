CREATE TABLE IF NOT EXISTS roles (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    level INT NOT NULL,
    description VARCHAR(255) NOT NULL
);


INSERT INTO roles(name , level , description) VALUES (
    'user',
    1,
    'Basic user'
);

INSERT INTO roles(name , level , description) VALUES (
    'admin',
    3,
    'Admin user'
);

INSERT INTO roles(name , level , description) VALUES (
    'mentor',
    2,
    'Mentor user'
);