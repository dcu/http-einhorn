package einhorn_test

import (
	"github.com/dcu/http-einhorn"
	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
	"os"
)

func (s *S) Test_EinhornNotRunning(c *C) {
	os.Setenv("EINHORN_FD_COUNT", "")
	c.Assert(einhorn.IsRunning(), Equals, false)
}

func (s *S) Test_EinhornRunning(c *C) {
	os.Setenv("EINHORN_FD_COUNT", "1")
	c.Assert(einhorn.IsRunning(), Equals, true)
}

func (s *S) Test_StartWithoutEinhorn(c *C) {
	g := gin.Default()
	err := einhorn.Start(g, 0)

	c.Assert(err, ErrorMatches, ".*too few EINHORN_FDs passed")
}

func (s *S) Test_StartWhenEinhornSocketIsClosed(c *C) {
	os.Setenv("EINHORN_FD_COUNT", "1")
	os.Setenv("EINHORN_FD_0", "1")
	g := gin.Default()
	err := einhorn.Start(g, 0)
	c.Assert(err, ErrorMatches, ".*socket operation on non-socket")
}
