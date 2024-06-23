package core

type Task struct {
	ID       string
	Commands []string
	Packages []string
    Variables map[string]string
}


