package lats_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-golang/localip"
	latsConfig "lats/config"
	"lats/helpers"
	"net/http"
	"os/exec"
	"testing"
)

var config *latsConfig.TestConfig

func TestLats(t *testing.T) {
	RegisterFailHandler(Fail)

	var metronSession *gexec.Session
	config = latsConfig.Load()

	BeforeSuite(func() {
		config.SaveMetronConfig()
		helpers.Initialize(config)
		metronSession = setupMetron()
	})

	AfterSuite(func() {
		metronSession.Kill().Wait()
	})

	RunSpecs(t, "Lats Suite")
}

func setupMetron() *gexec.Session {
	pathToMetronExecutable, err := gexec.Build("metron")
	Expect(err).ShouldNot(HaveOccurred())

	command := exec.Command(pathToMetronExecutable, "--config=fixtures/metron.json", "--debug")
	metronSession, err := gexec.Start(command, gexec.NewPrefixedWriter("[o][metron]", GinkgoWriter), gexec.NewPrefixedWriter("[e][metron]", GinkgoWriter))
	Expect(err).ShouldNot(HaveOccurred())

	localIPAddress, _ := localip.LocalIP()

	// wait for server to be up
	Eventually(func() error {
		_, err := http.Get("http://" + localIPAddress + ":1234")
		return err
	}, 3).ShouldNot(HaveOccurred())

	return metronSession
}
