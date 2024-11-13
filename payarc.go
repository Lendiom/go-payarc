package payarc

import (
	"encoding/json"
	"strings"
	"time"
)

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

// DateTime is a special handler for PayArc's various date formats
type DateTime struct {
	time.Time
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(time.RFC3339))
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	format := time.DateOnly
	if strings.Contains(str, "T") {
		format = time.RFC3339
	} else if strings.Contains(str, " ") {
		format = time.DateTime
	}

	parsed, err := time.Parse(format, str)
	if err != nil {
		return err
	}

	d.Time = parsed

	return nil
}

func (d DateTime) String() string {
	return d.Time.Format(time.DateTime)
}

func (d *DateTime) UnmarshalText(data []byte) error {
	str := string(data)

	format := time.DateOnly
	if strings.Contains(str, "T") {
		format = time.RFC3339
	} else if strings.Contains(str, " ") {
		format = time.DateTime
	}

	parsed, err := time.Parse(format, str)
	if err != nil {
		return err
	}

	d.Time = parsed

	return nil
}
