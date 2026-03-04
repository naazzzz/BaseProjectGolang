-- create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "username" character varying(191) NOT NULL,
  "password" text NOT NULL,
  "active" boolean NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- create "oauth_access_tokens" table
CREATE TABLE "oauth_access_tokens" (
  "id" text NOT NULL,
  "user_id" bigint NULL,
  "client_id" bigint NULL,
  "name" text NULL,
  "scopes" text NULL,
  "revoked" boolean NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "expires_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_oauth_access_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
