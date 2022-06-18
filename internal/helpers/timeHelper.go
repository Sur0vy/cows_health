package helpers

import (
	"math"
)

//2018-04-09T23:00:00Z
//const ctLayout = "2006-01-02|15:04:05"

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

//func ParseDateTime(input string) (time.Time, error) {
//	//example input (): "2014/08/01|11:27:18"
//	return time.Parse(ctLayout, input)
//}

//
//type CustomTime struct {
//	time.Time
//}
//
//func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
//	s := strings.Trim(string(b), "\"")
//	if s == "null" {
//		ct.Time = time.Time{}
//		return
//	}
//	ct.Time, err = time.Parse(ctLayout, s)
//	return
//}
//
//func (ct *CustomTime) MarshalJSON() ([]byte, error) {
//	if ct.Time.UnixNano() == nilTime {
//		return []byte("null"), nil
//	}
//	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
//}
//
//var nilTime = (time.Time{}).UnixNano()
//
//func (ct *CustomTime) IsSet() bool {
//	return ct.UnixNano() != nilTime
//}

//type Args struct {
//	Time CustomTime
//}
//
//var data = `
//    {"Time": "2014/08/01|11:27:18"}
//`

//func main() {
//	a := Args{}
//	fmt.Println(json.Unmarshal([]byte(data), &a))
//	fmt.Println(a.Time.String())
//}
