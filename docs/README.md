## Requirements

Your task will be to build a simplified version of the Reddit feed API that powers
http://old.reddit.com. **Please use Go as the programming language**. This will be a REST API
that allows users to do the following:

**Create new posts.** Posts should be validated for correctness. They should have the following
fields:

1. Title
2. Author: This should be a random 8 character string prefixed with t2_. The 8 character
string should only contain lowercase letters and numbers. For example, my user ID is
t2_11qnzrqv.
3. Link: This should be a valid URL. It's ok if your validation is not perfect.
4. Subreddit: The subreddit associated with this post.
5. Content: In the case of a text-only post. **A post cannot have both a link and content populated.**
6. Score: The total score associated with the upvotes and downvotes of a post.
7. Promoted: A boolean field indicating whether or not the post is an ad or not.
8. NSFW: Not safe for work. A boolean that indicates whether or not the post is safe for
work

**Generate a feed of posts.** This feed should have the following characteristics:

1. It should be ranked **by score**, and the post with the highest score should show up first.
2. It should be paginated, and each page should have **at most 27 posts**. Your API should
support fetching a specific page in the feed.
3. If a page has **3 posts or greater**, the second post should always be a promoted post if a
promoted post is available, regardless of the score.
4. If a page has **greater than 16 posts**, the 16th post should always be a promoted post if a promoted post is available, regardless of the score.
5. As an exception to rules 3 and 4, **a promoted post should never be shown adjacent
to an NSFW post**. You can ignore rules 3 and 4 in this case.
