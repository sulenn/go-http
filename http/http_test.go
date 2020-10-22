package http

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/sulenn/go-http/core/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sulenn/go-http/core/types"
)

var (
	url  = "http://023.node.internetapi.cn:21030/SCIDE/SCManager"
	body = &types.RequestJSON{
		Action: "executeContract",
	}
	params = map[string]string{"action": "executeContract"}
	db     *sql.DB

	// database
	username     = "root"
	password     = "qiubing"
	databaseName = "GITHUB"
)

func connectDB(t *testing.T) {
	database, err := sql.Open("mysql",
		username+":"+password+"@tcp(127.0.0.1:3306)/"+databaseName+"?charset=utf8")
	if err != nil {
		t.Fatalf("connect DB failed, %v\n", err)
	}
	db = database
}

func TestMain(m *testing.M) {
	connectDB(&testing.T{})
	exitCode := m.Run()
	err := db.Close()
	if err != nil {
		log.Fatalf("close db failed: %v\n", err)
	}
	os.Exit(exitCode)
}

func Test_GetAndParse(t *testing.T) {
	params["contractID"] = "RepositoryDB0"
	params["operation"] = "getCommit"
	params["arg"] = "{\"commit_hash\":\"082c254068a6564bc03516d9415930f90672e8dd\",\"repo_id\": \"205803088\"}"

	bytes, err := Get(url, params, nil)
	if err != nil {
		t.Fatal(err)
	}
	responseJSON, err := utils.ParseHttpResponse(bytes)
	if err != nil {
		t.Fatalf("parse bytes from response failed: %+v\n", err)
	}
	log.Printf("all info of response body is: %+v", responseJSON)
}

func Test_Post(t *testing.T) {
	body.ContractID = "RepositoryDB0"
	body.Operation = "putCommit"
	body.Arg = "{\"commit_hash\":\"3f2dc2ec21876d31ffa06744f215b6a17c5d3bay\",\"repo_id\":\"1002608\",\"commit_diff\":\"<273409891@qq.com>\"}"
	bytes, err := Post(url, body, params, nil)
	responseJSON, err := utils.ParseHttpResponse(bytes)
	if err != nil {
		t.Fatalf("parse bytes from response failed: %+v\n", err)
	}
	log.Printf("all info of response body is: %+v", responseJSON)
}

func Test_Post_ForPutCommit(t *testing.T) {
	body.Operation = "putCommit"
	repoID := 205803088
	rows, err := db.Query("SELECT * FROM github_commit WHERE repo_id = ?", repoID)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	commit := types.Commit{}
	var commitDiff sql.NullString
	for rows.Next() {
		if err := rows.Scan(&commit.CommitHash, &commit.RepoID, &commit.Author, &commit.Email,
			&commit.Time, &commit.Message, &commitDiff); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(commit)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if commit.RepoID%2 == 0 {
			body.ContractID = "RepositoryDB0"
		} else {
			body.ContractID = "RepositoryDB1"
		}
		body.Arg = string(jsonBytes)
		log.Printf("data in request body is: %+v\n", body)
		response, err := Post(url, body, params, nil)
		if err != nil {
			t.Fatalf("post request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForGetCommit(t *testing.T) {
	params["operation"] = "getCommit"
	repoID := 205803088
	rows, err := db.Query("SELECT commit_hash FROM github_commit WHERE repo_id = ?", repoID)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	commit := struct {
		CommitHash string `json:"commit_hash"`
		RepoID     int64  `json:"repo_id"`
	}{}
	commit.RepoID = int64(repoID)
	for rows.Next() {
		if err := rows.Scan(&commit.CommitHash); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(commit)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if commit.RepoID%2 == 0 {
			params["contractID"] = "RepositoryDB0"
		} else {
			params["contractID"] = "RepositoryDB1"
		}
		params["arg"] = string(jsonBytes)
		log.Printf("params in url is: %+v\n", params)
		response, err := Get(url, params, nil)
		if err != nil {
			t.Fatalf("get request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForPutIssue(t *testing.T) {
	body.Operation = "putIssue"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, number, user_login, created_at, "+
		"updated_at, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issue := types.Issue{}
	issue.OwnerName = ownerLogin
	issue.RepoName = repo
	issue.RepoID = 205803088
	var flag int
	for rows.Next() {
		if err := rows.Scan(&issue.IssueID, &issue.IssueNumber, &issue.UserName, &issue.CreatedAT,
			&issue.UpdatedAT, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 0 {
			jsonBytes, err := json.Marshal(issue)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if issue.RepoID%2 == 0 {
				body.ContractID = "RepositoryDB0"
			} else {
				body.ContractID = "RepositoryDB1"
			}
			body.Arg = string(jsonBytes)
			log.Printf("data in request body is: %+v\n", body)
			response, err := Post(url, body, params, nil)
			if err != nil {
				t.Fatalf("post request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForGetIssue(t *testing.T) {
	params["operation"] = "getIssue"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issue := struct {
		IssueID string `json:"issue_id"`
		RepoID  int64  `json:"repo_id"`
	}{}
	issue.RepoID = 205803088
	var flag int64
	for rows.Next() {
		if err := rows.Scan(&issue.IssueID, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 0 {
			jsonBytes, err := json.Marshal(issue)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if issue.RepoID%2 == 0 {
				params["contractID"] = "RepositoryDB0"
			} else {
				params["contractID"] = "RepositoryDB1"
			}
			params["arg"] = string(jsonBytes)
			log.Printf("params in url is: %+v\n", params)
			response, err := Get(url, params, nil)
			if err != nil {
				t.Fatalf("get request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForGetIssueList(t *testing.T) {
	params["operation"] = "getIssueList"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issue := struct {
		IssueID string `json:"issue_id"`
		RepoID  int64  `json:"repo_id"`
	}{}
	issue.RepoID = 205803088
	var flag int64
	for rows.Next() {
		if err := rows.Scan(&issue.IssueID, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 0 {
			jsonBytes, err := json.Marshal(issue)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if issue.RepoID%2 == 0 {
				params["contractID"] = "RepositoryDB0"
			} else {
				params["contractID"] = "RepositoryDB1"
			}
			params["arg"] = string(jsonBytes)
			log.Printf("params in url is: %+v\n", params)
			response, err := Get(url, params, nil)
			if err != nil {
				t.Fatalf("get request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForPutPullRequest(t *testing.T) {
	body.Operation = "putPullRequest"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, number, user_login, created_at, "+
		"updated_at, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequest := types.PullRequest{}
	pullRequest.OwnerName = ownerLogin
	pullRequest.RepoName = repo
	pullRequest.RepoID = 205803088
	//pullRequest.Content = "qiubing"
	var flag int
	for rows.Next() {
		if err := rows.Scan(&pullRequest.PullRequestID, &pullRequest.PullRequestNumber,
			&pullRequest.UserName, &pullRequest.CreatedAT,
			&pullRequest.UpdatedAT, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 1 {
			jsonBytes, err := json.Marshal(pullRequest)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if pullRequest.RepoID%2 == 0 {
				body.ContractID = "RepositoryDB0"
			} else {
				body.ContractID = "RepositoryDB1"
			}
			body.Arg = string(jsonBytes)
			log.Printf("data in request body is: %+v\n", body)
			response, err := Post(url, body, params, nil)
			if err != nil {
				t.Fatalf("post request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForGetPullRequest(t *testing.T) {
	params["operation"] = "getPullRequest"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequest := struct {
		PullRequestID string `json:"pull_request_id"`
		RepoID        int64  `json:"repo_id"`
	}{}
	pullRequest.RepoID = 205803088
	var flag int64
	for rows.Next() {
		if err := rows.Scan(&pullRequest.PullRequestID, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 1 {
			jsonBytes, err := json.Marshal(pullRequest)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if pullRequest.RepoID%2 == 0 {
				params["contractID"] = "RepositoryDB0"
			} else {
				params["contractID"] = "RepositoryDB1"
			}
			params["arg"] = string(jsonBytes)
			log.Printf("params in url is: %+v\n", params)
			response, err := Get(url, params, nil)
			if err != nil {
				t.Fatalf("get request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForGetPullRequestList(t *testing.T) {
	params["operation"] = "getPullRequestList"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, flag FROM github_issue WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequest := struct {
		PullRequestID string `json:"pull_request_id"`
		RepoID        int64  `json:"repo_id"`
	}{}
	pullRequest.RepoID = 205803088
	var flag int64
	for rows.Next() {
		if err := rows.Scan(&pullRequest.PullRequestID, &flag); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		if flag == 1 {
			jsonBytes, err := json.Marshal(pullRequest)
			if err != nil {
				t.Fatalf("struct to json failed, %v\n", err)
			}
			if pullRequest.RepoID%2 == 0 {
				params["contractID"] = "RepositoryDB0"
			} else {
				params["contractID"] = "RepositoryDB1"
			}
			params["arg"] = string(jsonBytes)
			log.Printf("params in url is: %+v\n", params)
			response, err := Get(url, params, nil)
			if err != nil {
				t.Fatalf("get request failed, %v\n", err)
			}
			log.Println(string(response))
		}
	}
}

func Test_Post_ForPutPullRequestComment(t *testing.T) {
	body.Operation = "putPullRequestComment"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, issue_number, user_login, created_at, "+
		"updated_at FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequestComment := types.PullRequestComment{}
	pullRequestComment.OwnerName = ownerLogin
	pullRequestComment.RepoName = repo
	pullRequestComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&pullRequestComment.PullRequestCommentID, &pullRequestComment.PullRequestNumber,
			&pullRequestComment.UserName, &pullRequestComment.CreatedAT,
			&pullRequestComment.UpdatedAT); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(pullRequestComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if pullRequestComment.RepoID%2 == 0 {
			body.ContractID = "RepositoryDB0"
		} else {
			body.ContractID = "RepositoryDB1"
		}
		body.Arg = string(jsonBytes)
		log.Printf("data in request body is: %+v\n", body)
		response, err := Post(url, body, params, nil)
		if err != nil {
			t.Fatalf("post request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForGetPullRequestComment(t *testing.T) {
	params["operation"] = "getPullRequestComment"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequestComment := struct {
		PullRequestCommentID string `json:"pull_request_comment_id"`
		RepoID               int64  `json:"repo_id"`
	}{}
	pullRequestComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&pullRequestComment.PullRequestCommentID); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(pullRequestComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if pullRequestComment.RepoID%2 == 0 {
			params["contractID"] = "RepositoryDB0"
		} else {
			params["contractID"] = "RepositoryDB1"
		}
		params["arg"] = string(jsonBytes)
		log.Printf("params in url is: %+v\n", params)
		response, err := Get(url, params, nil)
		if err != nil {
			t.Fatalf("get request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForGetPullRequestCommentList(t *testing.T) {
	params["operation"] = "getPullRequestCommentList"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	pullRequestComment := struct {
		PullRequestCommentID string `json:"pull_request_comment_id"`
		RepoID               int64  `json:"repo_id"`
	}{}
	pullRequestComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&pullRequestComment.PullRequestCommentID); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(pullRequestComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if pullRequestComment.RepoID%2 == 0 {
			params["contractID"] = "RepositoryDB0"
		} else {
			params["contractID"] = "RepositoryDB1"
		}
		params["arg"] = string(jsonBytes)
		log.Printf("params in url is: %+v\n", params)
		response, err := Get(url, params, nil)
		if err != nil {
			t.Fatalf("get request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForPutIssueComment(t *testing.T) {
	body.Operation = "putIssueComment"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id, issue_number, user_login, created_at, "+
		"updated_at FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issueComment := types.IssueComment{}
	issueComment.OwnerName = ownerLogin
	issueComment.RepoName = repo
	issueComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&issueComment.IssueCommentId, &issueComment.IssueNumber,
			&issueComment.UserName, &issueComment.CreatedAT,
			&issueComment.UpdatedAT); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(issueComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if issueComment.RepoID%2 == 0 {
			body.ContractID = "RepositoryDB0"
		} else {
			body.ContractID = "RepositoryDB1"
		}
		body.Arg = string(jsonBytes)
		log.Printf("data in request body is: %+v\n", body)
		response, err := Post(url, body, params, nil)
		if err != nil {
			t.Fatalf("post request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForGetIssueComment(t *testing.T) {
	params["operation"] = "getIssueComment"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issueComment := struct {
		IssueCommentID string `json:"issue_comment_id"`
		RepoID         int64  `json:"repo_id"`
	}{}
	issueComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&issueComment.IssueCommentID); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(issueComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if issueComment.RepoID%2 == 0 {
			params["contractID"] = "RepositoryDB0"
		} else {
			params["contractID"] = "RepositoryDB1"
		}
		params["arg"] = string(jsonBytes)
		log.Printf("params in url is: %+v\n", params)
		response, err := Get(url, params, nil)
		if err != nil {
			t.Fatalf("get request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}

func Test_Post_ForGetIssueCommentList(t *testing.T) {
	params["operation"] = "getIssueCommentList"
	ownerLogin := "FISCO-BCOS"
	repo := "go-sdk"
	rows, err := db.Query("SELECT id FROM github_comment WHERE owner_login=? and repo=?", ownerLogin, repo)
	if err != nil {
		t.Fatalf("query data failed from DB, %v\n", err)
	}
	issueComment := struct {
		IssueCommentID string `json:"issue_comment_id"`
		RepoID         int64  `json:"repo_id"`
	}{}
	issueComment.RepoID = 205803088
	for rows.Next() {
		if err := rows.Scan(&issueComment.IssueCommentID); err != nil {
			t.Fatalf("read data failed from rows.Scan, %v\n", err)
		}
		jsonBytes, err := json.Marshal(issueComment)
		if err != nil {
			t.Fatalf("struct to json failed, %v\n", err)
		}
		if issueComment.RepoID%2 == 0 {
			params["contractID"] = "RepositoryDB0"
		} else {
			params["contractID"] = "RepositoryDB1"
		}
		params["arg"] = string(jsonBytes)
		log.Printf("params in url is: %+v\n", params)
		response, err := Get(url, params, nil)
		if err != nil {
			t.Fatalf("get request failed, %v\n", err)
		}
		log.Println(string(response))
	}
}
