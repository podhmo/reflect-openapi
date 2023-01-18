package jsonequal

import (
	"fmt"

	"github.com/nsf/jsondiff"
)

// FailPlain :
func FailPlain(
	c *Caller,
	left interface{},
	right interface{},
	lb []byte,
	rb []byte,
) error {
	ls, rs := string(lb), string(rb)
	if ls == rs {
		msg := "%[5]s%[6]s (%[1]T):\n	%[3]s\n%[7]s (%[2]T):\n	%[4]s"
		return fmt.Errorf(msg, left, right, ls, rs, c.Prefix, c.LeftName, c.RightName)
	}
	// todo : more redable expression
	msg := "%[5]s%[6]s:\n	%[3]s\n%[7]s:\n	%[4]s"
	return fmt.Errorf(msg, left, right, ls, rs, c.Prefix, c.LeftName, c.RightName)
}

// FailJSONDiff :
func FailJSONDiff(
	c *Caller,
	left interface{},
	right interface{},
	lb []byte,
	rb []byte,
) error {
	options := jsondiff.DefaultConsoleOptions()
	diff, s := jsondiff.Compare(lb, rb, &options)
	prefix := c.Prefix
	c.Prefix = "\n"
	return fmt.Errorf("%sstatus=%s\n%s\n%s", prefix, diff.String(), s, FailPlain(c, left, right, lb, rb).Error())
}
