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

CREATE INDEX score_desc_index ON post (score DESC NULLS LAST);
CREATE INDEX promoted_post_index ON post USING btree(promoted) WHERE promoted is TRUE;

-- +goose Down
DROP INDEX promoted_post_index;
DROP INDEX score_desc_index;

DROP TABLE post;
