CREATE TABLE "users" (
  "telegramID"    INTEGER     NOT NULL    PRIMARY KEY,
  id              SERIAL      NOT NULL,
  anonymous       BOOLEAN     NOT NULL,
  admin           BOOLEAN     NOT NULL    DEFAULT FALSE
);

CREATE UNIQUE INDEX index_on_users_id
  ON "users"(id);

CREATE TYPE ordering AS ENUM ('first', 'last', 'random');
CREATE TABLE questions (
  id              SERIAL      NOT NULL    PRIMARY KEY,
  body            TEXT        NOT NULL,
  "order"         ordering    NOT NULL    DEFAULT 'random'
);

CREATE TYPE status AS ENUM ('requested', 'asked', 'done', 'failed');
CREATE TABLE pauses (
  id              SERIAL      NOT NULL    PRIMARY KEY,
  "user"          INT         NOT NULL    REFERENCES "users"(id),
  question        INT         NOT NULL    REFERENCES questions(id),
  answer          TEXT,
  status          status      NOT NULL    DEFAULT 'requested',
  message_id      INT,
  chat_id         BIGINT
);

-- CREATE TABLE mailings (
--   date DATETIME
-- )
