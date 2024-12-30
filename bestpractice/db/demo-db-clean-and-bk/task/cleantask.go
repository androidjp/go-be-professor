package task

import (
	"demodbclient/repo"
	"fmt"

	"github.com/robfig/cron/v3"
)

// CleanTask 定期清理任务
type CleanTask struct {
	cr *cron.Cron
	r  *repo.Repo
}

func NewCleanTask(r *repo.Repo) *CleanTask {
	return &CleanTask{
		cr: cron.New(
			cron.WithParser(cron.NewParser(
				cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
			)),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			)),
		r: r}
}

func (c *CleanTask) Run() {
	taskID, err := c.cr.AddFunc("*/5 * * * * *", c.Clean)
	if err != nil {
		fmt.Printf("失败，err：%s", err.Error())
		return
	}
	c.cr.Start()
	fmt.Printf("taskID: %d", taskID)
}

func (c *CleanTask) Clean() {

	// DB清理动作
	c.r.DeleteOutdated()

	fmt.Println("clean..........")
}
