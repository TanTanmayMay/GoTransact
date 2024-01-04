--liquibase formatted sql

--changeset tanmay-kokate:1

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--rollback DROP EXTENSION "uuid-ossp";

--changeset tanmay-kokate:2

CREATE TABLE users ("id" UUID, "name" VARCHAR(50) NOT NULL, "password" VARCHAR(50) NOT NULL, PRIMARY KEY("id"));


--rollback DROP TABLE workspace;