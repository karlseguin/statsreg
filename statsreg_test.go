package statsreg

import (
	. "github.com/karlseguin/expect"
	"testing"
	"io/ioutil"
	"encoding/json"
)

type StatsRegTest struct{}

func Test_StatsReg(t *testing.T) {
	Expectify(new(StatsRegTest), t)
}

func (_ StatsRegTest) CollectsAndOutputStats() {
	sr := New(Configure().File("test.json", true))
	sr.RegisterString("last", ProvideString)
	sr.RegisterInt64("power", ProvideInt64)
	sr.Collect()
	assertFile("last", "it's over", "power", float64(9000))
}

func ProvideString() string {
	return "it's over"
}

func ProvideInt64() int64 {
	return 9000
}

func assertFile(keyValues ...interface{}) {
	bytes, err :=	ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}
	var stats map[string]interface{}
	if err := json.Unmarshal(bytes, &stats); err != nil {
		panic(err)
	}

	Expect(len(stats)).To.Equal(len(keyValues)/2)
	for i := 0; i < len(keyValues); i += 2 {
		key, value := keyValues[i].(string), keyValues[i+1]
		Expect(stats[key]).To.Equal(value)
	}
}
