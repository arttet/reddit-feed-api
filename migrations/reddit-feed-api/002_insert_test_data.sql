-- +goose Up

INSERT INTO post
  (title, author, link, subreddit, content, score, promoted, not_safe_for_work)
VALUES
  ('Title', 't2_00000001', 'https://reddit.com', 'SubReddit', '', 85, false, false),
  ('Title', 't2_00000002', '', 'SubReddit', 'Content', 86, false, false),
  ('Title', 't2_00000003', 'https://reddit.com', 'SubReddit', '', 87, false, false),
  ('Title', 't2_00000004', '', 'SubReddit', 'Content', 88, false, false),
  ('Title', 't2_00000005', 'https://reddit.com', 'SubReddit', '', 89, false, false),
  ('Title', 't2_00000006', '', 'SubReddit', 'Content', 90, false, false),
  ('Title', 't2_00000007', 'https://reddit.com', 'SubReddit', '', 91, false, false),
  ('Title', 't2_00000008', '', 'SubReddit', 'Content', 92, false, false),
  ('Title', 't2_00000009', 'https://reddit.com', 'SubReddit', '', 93, false, false),
  ('Title', 't2_00000010', '', 'SubReddit', 'Content', 94, false, false),
  ('Title', 't2_00000011', 'https://reddit.com', 'SubReddit', '', 95, false, false),
  ('Title', 't2_00000012', '', 'SubReddit', 'Content', 96, false, false),
  ('Title', 't2_00000013', 'https://reddit.com', 'SubReddit', '', 97, false, false),
  ('Title', 't2_00000014', '', 'SubReddit', 'Content', 98, false, false),
  ('Title', 't2_00000015', 'https://reddit.com', 'SubReddit', '', 99, false, false),
  ('Title', 't2_00000016', '', 'SubReddit', 'Content', 100, false, true),
  ('Title', 't2_00000017', 'https://reddit.com', 'SubReddit', '', 0, true, false);

-- +goose Down

DELETE FROM post;
