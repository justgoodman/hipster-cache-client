package test

import (
	. "gopkg.in/check.v1"
	"testing"
)

type StringTestSuite struct{}

var _ = Suite(&StringTestSuite{})

func TestStart(t *testing.T) {
	TestingT(t)
}

func (suite *StringTestSuite) TestSetExpectNotError(c *C) {
	client := NewClientDSL().Do()

	err := client.Set("test", "value")

	c.Assert(err, IsNil)
}

func (suite *StringTestSuite) TestSetEmptyKeyExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.Set("", "value")

	c.Assert(err, NotNil)
}

func (suite *StringTestSuite) TestSetKeyMoreThanMaxLenExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.Set(RandString(256), "value")

	c.Assert(err, NotNil)

}
func (suite *StringTestSuite) TestGetExpectObservedValue(c *C) {
	key := "some_key"
	value := "some_value"
	client := NewClientDSL().Do()
	client.Set(key, value)

	retValue, err := client.Get(key)

	c.Assert(err, IsNil)
	c.Assert(retValue, Equals, value)
}

func (suite *StringTestSuite) TestGetEmptyKeyExpectError(c *C) {
	client := NewClientDSL().Do()

	_, err := client.Get("")

	c.Assert(err, NotNil)
}

func (suite *StringTestSuite) TestGetKeyMoreThanMaxLenExpectError(c *C) {
	client := NewClientDSL().Do()

	_, err := client.Get(RandString(256))

	c.Assert(err, NotNil)
}
