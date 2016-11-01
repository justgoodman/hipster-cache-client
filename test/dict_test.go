package test

import (
	. "gopkg.in/check.v1"
)

type DictTestSuite struct{}

var _ = Suite(&DictTestSuite{})

func (suite *DictTestSuite) TestDSetExpectNotError(c *C) {
	client := NewClientDSL().Do()

	err := client.DSet("testDictKey", "field", "value")

	c.Assert(err, IsNil)
}

func (suite *DictTestSuite) TestDSetEmptyFielValueError(c *C) {
	client := NewClientDSL().Do()

	err := client.DSet("testDictKey", "", "")

	c.Assert(err, NotNil)
}

func (suite *DictTestSuite) TestDSetEmptyKeyExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.DSet("", "field", "value")

	c.Assert(err, NotNil)
}

func (suite *DictTestSuite) TestDSetMoreThanMaxLenExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.DSet(RandString(256), "field", "value")

	c.Assert(err, NotNil)

}

func (suite *DictTestSuite) TestDGetExpectObservedValue(c *C) {
	key := "testDictKey"
	field := "field"
	value := "value"
	client := NewClientDSL().Do()
	client.DSet(key, field, value)

	retValue, err := client.DGet(key, field)

	c.Assert(err, IsNil)
	c.Assert(retValue, Equals, value)
}
