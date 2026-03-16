-- create "example_models" table
CREATE TABLE "example_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "data" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_example_models_deleted_at" to table: "example_models"
CREATE INDEX "idx_example_models_deleted_at" ON "example_models" ("deleted_at");
