package style

const DefaultMaxOpts = 5

var DefaultDistribution = HorizontalDistribution(JustifyAround, DefaultMaxOpts)

type Distribution struct {
	Direction Direction
	Justify   Justify
	Limit     uint16
}

func VerticalDistribution() Distribution {
	return Distribution{
		Direction: Vertical,
	}
}

func HorizontalDistribution(justify Justify, limit uint16) Distribution {
	return Distribution{
		Direction: Horizontal,
		Justify:   justify,
		Limit:     limit,
	}
}
