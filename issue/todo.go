package issue

type todo struct {
	Id      uint64 `json:"todo_id"`
	Index   uint32 `json:"index"`
	IssueId uint64 `json:"issue_id"`
	Active  bool   `json:"active"`
}
