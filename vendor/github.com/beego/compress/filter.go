package compress

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var (
	ClosureBin  = "java -jar compiler.jar"
	ClosureArgs = map[string]string{
		"compilation_level": "SIMPLE_OPTIMIZATIONS",
		"warning_level":     "QUIET",
	}

	YuiBin  = "java -jar yuicompressor.jar"
	YuiArgs = map[string]string{
		"type": "css",
	}
)

func ClosureFilter(source string) string {
	args := strings.Fields(ClosureBin)
	for arg, value := range ClosureArgs {
		args = append(args, "--"+arg)
		args = append(args, value)
	}
	return runFilter(args[0], args[1:], source)
}

func YuiFilter(source string) string {
	args := strings.Fields(YuiBin)
	for arg, value := range YuiArgs {
		args = append(args, "--"+arg)
		args = append(args, value)
	}
	return runFilter(args[0], args[1:], source)
}

func runFilter(bin string, args []string, source string) string {
	buf := bytes.NewBufferString(source)
	out := bytes.NewBufferString("")

	cmd := exec.Command(bin, args...)
	cmd.Stdin = buf
	cmd.Stderr = os.Stderr
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		logError(err.Error())
		return source
	} else {
		return out.String()
	}
}
