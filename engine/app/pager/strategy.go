package pager

var default_engine = EnginePage()
var default_predicate = PredicatePage()

type PagerStrategy struct {
	Engine    Engine
	Predicate Predicate
}

func NewStrategy() PagerStrategy {
	return PagerStrategy{
		Engine:    default_engine,
		Predicate: default_predicate,
	}
}

func (p *PagerStrategy) SetEngine(engine Engine) *PagerStrategy {
	p.Engine = engine
	return p
}

func (p *PagerStrategy) SetPredicate(predicate Predicate) *PagerStrategy {
	p.Predicate = predicate
	return p
}
