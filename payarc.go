package payarc

type Boolean uint8

var (
	False Boolean = 0
	True  Boolean = 1
)

func (b Boolean) AsBool() bool {
	return b == True
}

type YesOrNo string

var (
	Yes YesOrNo = "yes"
	No  YesOrNo = "no"
)

func (y YesOrNo) AsBool() bool {
	return y == Yes
}

type ChargeCardLevel string

var (
	ChargeCardLevel1 ChargeCardLevel = "LEVEL1"
	ChargeCardLevel2 ChargeCardLevel = "LEVEL2"
	ChargeCardLevel3 ChargeCardLevel = "LEVEL3"
)
