CREATE TABLE "users" (
  "telegramID"    INTEGER     NOT NULL    PRIMARY KEY,
  id              SERIAL      NOT NULL,
  anonymous       BOOLEAN     NOT NULL
);

CREATE UNIQUE INDEX index_on_users_id
  ON "users"(id);

CREATE TYPE ordering AS ENUM ('first', 'last', 'random');

CREATE TABLE questions (
  id              SERIAL      NOT NULL    PRIMARY KEY,
  body            TEXT        NOT NULL,
  "order"         ordering    NOT NULL    DEFAULT 'random'
);

CREATE TABLE answers (
  id              SERIAL      NOT NULL    PRIMARY KEY,
  body            TEXT        NOT NULL,
  question        INT         NOT NULL    REFERENCES questions(id),
  "user"          INT         NOT NULL    REFERENCES "users"(id)
);

-- CREATE TABLE mailings (
--   date DATETIME
-- )