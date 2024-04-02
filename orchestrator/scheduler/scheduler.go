package scheduler

type Scheduler interface {
	SelectCandidateNode() string
	Score() int
	Pick() string
}
