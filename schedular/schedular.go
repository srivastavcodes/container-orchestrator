package schedular

type Schedular interface {
	SelectCandidateNodes()
	Score()
	Pick()
}
