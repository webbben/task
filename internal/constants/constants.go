package constants

type taskStatuses struct {
	Pending    string
	InProgress string
	Complete   string
}

var TaskStatus taskStatuses = taskStatuses{
	Pending:    "pending",
	InProgress: "in progress",
	Complete:   "complete",
}
