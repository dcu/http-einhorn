package einhorn_test

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
}

func (s *S) TearDownSuite(c *C) {
}

func (s *S) SetUpTest(c *C) {
}

func (s *S) TearDownTest(c *C) {
	os.Setenv("EINHORN_FD_COUNT", "")
}
