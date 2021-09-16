-- +goose Up
CREATE TABLE post (
  id                BIGSERIAL     PRIMARY KEY,
  title             TEXT          NOT NULL,
  author            CHARACTER(11) NOT NULL,
  link              TEXT          NOT NULL,
  subreddit         TEXT          NOT NULL,
  content           TEXT          NOT NULL,
  score             BIGINT        NOT NULL,
  promoted          BOOLEAN       NOT NULL,
  not_safe_for_work BOOLEAN       NOT NULL
);

-- +goose Down
DROP TABLE post;
