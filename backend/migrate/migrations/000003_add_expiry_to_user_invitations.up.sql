ALTER TABLE user_invitations
    ADD COLUMN expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL
    DEFAULT (NOW() + interval '50 minutes');