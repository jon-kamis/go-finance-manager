package logtest

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jon-kamis/klogger"
)

func GetTestPropertyFileName() string {
	fp, err := os.Getwd()

	if err != nil {
		panic("unexpected error when getting current directory")
	}

	//Do this to handle both windows and linux file systems

	r1 := regexp.MustCompile(`internal\\.*`)
	r2 := regexp.MustCompile(`internal/.*`)

	fp = r1.ReplaceAllString(fp, "")
	fp = r2.ReplaceAllString(fp, "")

	fn := "properties\\klogger-properties-test.yml"
	ffn := fmt.Sprintf("%s%s", fp, fn)

	fmt.Printf("Klogger Property file: %s\n", ffn)

	return ffn
}

func SetKloggerTestFileNameEnv() {
	os.Setenv("KloggerPropFileName", GetTestPropertyFileName())
	klogger.RefreshConfig()
}
