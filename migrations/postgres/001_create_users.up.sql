CREATE TYPE "user_roles" AS ENUM (
    'SYSTEM_ADMIN',
    'ORGANIZATION_ADMIN',
    'INSTRUCTOR',
    'STUDENT'
);

CREATE TABLE IF NOT EXISTS "user" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "role" user_roles NOT NULL,
    "name" VARCHAR NOT NULL,
    "phone" VARCHAR,
    "email" VARCHAR NOT NULL UNIQUE,
    "password" VARCHAR(1000),
    "active" SMALLINT,
    "expires_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON "user"(email);
