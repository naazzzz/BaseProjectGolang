-- create "users" table
CREATE TABLE `users` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `username` varchar NOT NULL,
  `password` text NOT NULL,
  `active` numeric NULL
);
-- create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX `idx_users_username` ON `users` (`username`);
-- create "oauth_access_tokens" table
CREATE TABLE `oauth_access_tokens` (
  `id` text NULL,
  `user_id` integer NULL,
  `client_id` integer NULL,
  `name` text NULL,
  `scopes` text NULL,
  `revoked` numeric NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `expires_at` datetime NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_oauth_access_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
