package issue

import (
	"fmt"
	"time"
)

type issue struct {
	Id           uint64     `json:"issue_id"`
	Label        string     `json:"label"`
	Description  string     `json:"description"`
	CreatedTime  time.Time  `json:"created_time"`
	UpdatedTime  *time.Time `json:"updated_time"`
	FinishedTime *time.Time `json:"finishied_time"`
	Priority     bool       `json:"priority"`
}

func main() {
	fmt.Println("vim-go")
}
