package style

const DefaultLimit = 5

var DefaultDistribution = HorizontalDistribution(JustifyAround, DefaultLimit)

type Distribution struct {
	Direction Direction
	Justify   Justify
	Limit     uint8
}

func VerticalDistribution() Distribution {
	return Distribution{
		Direction: Vertical,
	}
}

func HorizontalDistribution(justify Justify, limit uint8) Distribution {
	return Distribution{
		Direction: Horizontal,
		Justify:   justify,
		Limit:     limit,
	}
}
