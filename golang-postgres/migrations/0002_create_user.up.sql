CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    Name TEXT NOT NULL,
    Age BIGINT NOT NULL,
    city TEXT NOT NULL
);
