package objects

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Probe tests", func() {
	var (
		cmd     []string
		timeout int32
	)
	BeforeEach(func() {
		cmd = []string{"echo", "test"}
		timeout = 3
	})
	It("generate exec probe", func() {
		probe := getExecProbe(cmd, timeout)
		Expect(probe.Exec).ToNot(BeNil())
		Expect(probe.Exec.Command).Should(Equal(cmd))
		Expect(probe.TimeoutSeconds).Should(Equal(timeout))
	})
})
