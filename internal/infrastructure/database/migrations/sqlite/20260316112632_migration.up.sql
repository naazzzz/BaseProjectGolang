-- create "example_models" table
CREATE TABLE `example_models` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `data` text NOT NULL
);
-- create index "idx_example_models_deleted_at" to table: "example_models"
CREATE INDEX `idx_example_models_deleted_at` ON `example_models` (`deleted_at`);
