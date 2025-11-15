
CREATE TABLE IF NOT EXISTS user_invitation (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    token bytea NOT NULL,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT now()
) 

