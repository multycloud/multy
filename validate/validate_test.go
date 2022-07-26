package validate_test

import (
	"github.com/multycloud/multy/validate"
	"testing"
)

func testRegexp(validator *validate.RegexpValidator, shouldMatch, shouldntMatch []string, t *testing.T) {
	for _, name := range shouldMatch {
		if err := validator.Check(name, "some_val"); err != nil {
			t.Errorf("%v should match %s, but didn't", validator, name)
		}
	}
	for _, name := range shouldntMatch {
		if err := validator.Check(name, "some_val"); err == nil {
			t.Errorf("%v shouldn't match %s, but did", validator, name)
		}
	}
}

// TestWordWithDotHyphenUnder80Pattern checks whether validate.wordWithDotHyphenUnder80Pattern matches
// expected expressions
func TestWordWithDotHyphenUnder80Pattern(t *testing.T) {
	shouldMatch := []string{
		"a",
		"9",
		"aZ9",
		"ThisIs67dots..................................................................._",
	}
	shouldntMatch := []string{
		"",
		"_",
		"<someName",
		"ThisIs68dots...................................................................._",
		"Maybe?inThe.Middle_",
	}
	testRegexp(validate.NewWordWithDotHyphenUnder80Validator(), shouldMatch, shouldntMatch, t)
}

// TestCIDRIPv4Matching checks the correctness of validate.cidrIPv4Pattern
func TestCIDRIPv4Matching(t *testing.T) {
	shouldMatch := []string{
		"10.0.0.1",
		"0.0.0.0/0",
		"255.255.255.255/32",
		"172.16.0.0/16",
	}
	shouldntMatch := []string{
		"This is not CIDR",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	}
	testRegexp(validate.NewCIDRIPv4Check(), shouldMatch, shouldntMatch, t)
}

// TestProtocolMatching checks if only allowed protocol values match
func TestProtocolMatching(t *testing.T) {
	shouldMatch := []string{
		"tcp",
		"udp",
		"icmp",
		"*",
	}
	shouldntMatch := []string{
		"TCP",
		"ah",
		"IcmpV6",
		"ESP",
		"Oh",
		"*anything",
	}
	testRegexp(validate.NewProtocolCheck(), shouldMatch, shouldntMatch, t)
}

func TestPortRangeCheck(t *testing.T) {
	portCheck := validate.NewPortCheck()
	ok := []int32{0, 80, 8080, 443, 22, 65535}
	notOk := []int32{-1, 65536}
	for _, v := range ok {
		if err := portCheck.Check(v, "port"); err != nil {
			t.Errorf("%v should match, but didn't", v)
		}
	}
	for _, v := range notOk {
		if err := portCheck.Check(v, "port"); err == nil {
			t.Errorf("%v shouln't match, but did", v)
		}
	}
}

func TestPriorityCheck(t *testing.T) {
	priorityCheck := validate.NewPriorityCheck()
	ok := []int64{100, 4096, 101, 202}
	notOk := []int64{99, 4097, 0, -1}
	for _, v := range ok {
		if err := priorityCheck.Check(v, "priority"); err != nil {
			t.Errorf("%v should match, but didn't", v)
		}
	}
	for _, v := range notOk {
		if err := priorityCheck.Check(v, "priority"); err == nil {
			t.Errorf("%v shouln't match, but did", v)
		}
	}
}
