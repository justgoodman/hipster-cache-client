package test

import (
	. "gopkg.in/check.v1"
)

type ListTestSuite struct{}

var _ = Suite(&ListTestSuite{})

func (suite *ListTestSuite) TestLPushExpectNotError(c *C) {
	client := NewClientDSL().Do()

	err := client.LPush("testLPushKey","value")

	c.Assert(err, IsNil)
}

func (suite *ListTestSuite) TestLPushEmptyKeyExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.LPush("","value")

	c.Assert(err, NotNil)
}

func (suite *ListTestSuite) TestLPushMoreThanMaxLenExpectError(c *C) {
	client := NewClientDSL().Do()

	err := client.LPush(RandString(256), "value")

	c.Assert(err, NotNil)

}

func (suite *ListTestSuite) TestLLenExpectMoreThan0(c *C) {
	key := "testLPush"
	client := NewClientDSL().Do()
	client.LPush(key, "value")

	listLen, err := client.LLen(key)

	c.Assert(err, IsNil)
	c.Assert(listLen > 0, Equals, true)
}

func (suite *ListTestSuite) TestLRangeExpectObservedValues(c *C) {
	key := "testLPush"
	value1 := RandString(10)
	value2 := RandString(10)
	client := NewClientDSL().Do()
	startIndex,_ := client.LLen(key)
	client.LPush(key, value1)
	client.LPush(key, value2)

	values, err := client.LRange(key, startIndex, startIndex + 2)

	c.Assert(err, IsNil)
	c.Assert(values, DeepEquals, []string{value1,value2})
}

func (suite *ListTestSuite) TestLSetExpectObservedValue(c *C) {
	key := "testLPush"
	valuePush := RandString(10)
	valueSet := RandString(10)
	client := NewClientDSL().Do()
	client.LPush(key, valuePush)
	index,_ := client.LLen(key)

	err := client.LSet(key, index-1, valueSet)

	c.Assert(err, IsNil)

	value, _  := client.LRange(key, index-1, index-1)

	c.Assert(value, DeepEquals, []string{valueSet})
}
