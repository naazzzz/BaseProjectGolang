-- reverse: create index "idx_test_models_deleted_at" to table: "test_models"
DROP INDEX "idx_test_models_deleted_at";
-- reverse: create "test_models" table
DROP TABLE "test_models";
-- reverse: create index "idx_example_models_deleted_at" to table: "example_models"
DROP INDEX "idx_example_models_deleted_at";
-- reverse: create "example_models" table
DROP TABLE "example_models";
