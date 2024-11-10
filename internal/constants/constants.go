package constants

type taskStatuses struct {
	Pending    int
	InProgress int
	Complete   int
}

var TaskStatus taskStatuses = taskStatuses{
	Pending:    0,
	InProgress: 1,
	Complete:   10,
}

var TaskStatusDisplay = map[int]string{
	TaskStatus.Pending:    "waiting",
	TaskStatus.InProgress: "in prog",
	TaskStatus.Complete:   "COMP",
}
