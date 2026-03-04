#!/bin/sh
set -e

atlas migrate diff migration --env gorm_postgres

atlas migrate diff migration --env gorm_sqlite
