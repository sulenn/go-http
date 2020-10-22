package types

type Commit struct {
	CommitHash string `json:"commit_hash"`
	RepoID     int64  `json:"repo_id"`
	Author     string `json:"author"`
	Committer  string `json:"committer"`
	Email      string `json:"email"`
	Time       string `json:"time"`
	Message    string `json:"message"`
	CommitDiff string `json:"commit_diff"`
}

type PullRequest struct {
	PullRequestID     int64    `json:"pull_request_id"`
	PullRequestNumber int64    `json:"pull_request_number"`
	RepoID            int64    `json:"repo_id"`
	RepoName          string   `json:"reponame"`
	OwnerName         string   `json:"ownername"`
	UserName          string   `json:"username"`
	Action            string   `json:"action"`
	Title             string   `json:"title"`
	Content           string   `json:"content"`
	SourceBranch      string   `json:"source_branch"`
	TargetBranch      string   `json:"target_branch"`
	Reviewers         []string `json:"reviewers"`
	CommitSHAs        []string `json:"commit_shas"`
	MergeUser         string   `json:"merge_user"`
	CreatedAT         string   `json:"created_at"`
	UpdatedAT         string   `json:"updated_at"`
}

type Issue struct {
	IssueID     int64  `json:"issue_id"`
	IssueNumber int64  `json:"issue_number"`
	RepoID      int64  `json:"repo_id"`
	RepoName    string `json:"reponame"`
	OwnerName   string `json:"ownername"`
	UserName    string `json:"username"`
	Action      string `json:"action"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	CreatedAT   string `json:"created_at"`
	UpdatedAT   string `json:"updated_at"`
}

type CommentCommon struct {
	RepoID    int64  `json:"repo_id"`
	RepoName  string `json:"reponame"`
	OwnerName string `json:"ownername"`
	UserName  string `json:"username"`
	Action    string `json:"action"`
	Content   string `json:"content"`
	CreatedAT string `json:"created_at"`
	UpdatedAT string `json:"updated_at"`
}

type PullRequestComment struct {
	PullRequestCommentID int64 `json:"pull_request_comment_id"`
	PullRequestNumber    int64 `json:"pull_request_number"`
	CommentCommon
}

type IssueComment struct {
	IssueCommentId int64 `json:"issue_comment_id"`
	IssueNumber    int64 `json:"issue_number"`
	CommentCommon
}
