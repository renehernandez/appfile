package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	// "os"

	"testing"

	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
}

func (suite *RootTestSuite) TestMissingEnvFileLoadAtDefaultLocation() {
	cmd := rootCmd{}

	err := cmd.loadEnvVars()

	suite.Error(err)
}
func (suite *RootTestSuite) TestLoadsEnvFromFilePath() {
	envFile, err := createTempFile(map[string]string{"HELLO": "WORLD"})
	suite.NoError(err)
	defer os.Remove(envFile)

	cmd := rootCmd{
		envFile: envFile,
	}

	err = cmd.loadEnvVars()

	suite.NoError(err)
	suite.Equal("WORLD", os.Getenv("HELLO"))
}

func (suite *RootTestSuite) TestInitializeTokenSetFromDotenvFile() {
	envFile, err := createTempFile(map[string]string{"DIGITALOCEAN_ACCESS_TOKEN": "TOKEN"})
	suite.NoError(err)
	defer os.Remove(envFile)

	cmd := rootCmd{
		envFile:  envFile,
		logLevel: "debug",
	}

	err = cmd.initialize()

	suite.NoError(err)
	suite.Equal("TOKEN", os.Getenv("DIGITALOCEAN_ACCESS_TOKEN"))
}

func createTempFile(envData map[string]string) (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "env-")
	if err != nil {
		return "", err
	}

	for key, value := range envData {
		tmpFile.WriteString(fmt.Sprintf("%s=%s", key, value))
	}

	return tmpFile.Name(), nil
}

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, &RootTestSuite{})
}
