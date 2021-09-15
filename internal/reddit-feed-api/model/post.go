package model

type Post struct {
	ID             uint64 `db:"id"                yaml:"id"`
	Title          string `db:"title"             yaml:"title"`
	Author         string `db:"author"            yaml:"author"`
	Link           string `db:"link"              yaml:"link"`
	Subreddit      string `db:"subreddit"         yaml:"subreddit"`
	Content        string `db:"content"           yaml:"content"`
	Score          uint64 `db:"score"             yaml:"score"`
	Promoted       bool   `db:"promoted"          yaml:"promoted"`
	NotSafeForWork bool   `db:"not_safe_for_work" yaml:"notSafeForWork"`
}
