package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testFilePath = ".github/classroom/autograding.json"

type autogradingTest struct {
	Name   string
	Run    string
	Output string
	Points int
}

type autogradingTests struct {
	Tests []autogradingTest
}

var logString string
var score int

func main() {
	testsJson, err := os.ReadFile(testFilePath)
	if err != nil {
		fmt.Printf("something went wrong opening file: %s", testFilePath)
		os.Exit(1)
	}

	var tests autogradingTests
	err = json.Unmarshal(testsJson, &tests)
	if err != nil {
		fmt.Printf("something went wrong unmarshaling the tests: %s\n", err)
		os.Exit(1)
	}

	for _, test := range tests.Tests {
		// We need to check for sequenced commands
		var cmd []string
		if strings.Contains(test.Run, " && ") {
			cmds := strings.SplitN(test.Run, " && ", 2)
			cmd1 := strings.Split(cmds[0], " ")
			_, err := getCmdOutput(cmd1[0], cmd1[1:])
			if err != nil {
				fmt.Printf("something went wrong running command %s, error: %s\n", cmds[0], err)
				logString += fmt.Sprintf("something went wrong running command %s, error: %s\n", cmds[0], err)
			}
			cmd2 := cmds[1]
			cmd = strings.Split(cmd2, " ")
		} else {
			cmd1 := strings.Split(test.Run, " ")
			for _, c := range cmd1 {
				if strings.TrimSpace(c) != "" {
					c = strings.ReplaceAll(c, "\"", "")
					cmd = append(cmd, c)
				}
			}
		}
		logInfo("running: " + test.Name + " with command '" + test.Run + "' ...")
		out, err := getCmdOutput(cmd[0], cmd[1:])
		if err != nil {
			fmt.Printf("something went wrong running test command: %s, error: %s\n", test.Run, err)
			logString += fmt.Sprintf("something went wrong running test command: %s, error: %s\n", test.Run, err)
		}
		if out == test.Output {
			score += test.Points
			logSuccess("test " + test.Name + " succeeded! Score: " + fmt.Sprint(score+test.Points) + "/100")
		} else {
			logFailure("test "+test.Name+" failed. Output of command "+test.Run+" did not match expected. Score: "+fmt.Sprint(score)+"/100", test.Output, string(out))
		}
	}
	os.WriteFile("tester-log.txt", []byte(logString), os.FileMode(0777))
	fmt.Println("\n\nTests complete. Score: " + fmt.Sprintf("%d", score) + "/100. Test log written to tester-log.txt")
}

func getCmdOutput(cmdStr string, args []string) (string, error) {
	cmd := exec.Command(cmdStr, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("something went wrong getting command output, command %s, error: %s\n", cmd, err)
		return "", err
	}
	return string(out), nil
}

func logInfo(s string) {
	fmt.Println("tester: " + s)
	logString += "tester: " + s + "\n"
}

func logFailure(s string, expected string, actual string) {
	expected, actual = addTabs(expected), addTabs(actual)
	fmt.Println("tester: FAILURE: " + s)
	fmt.Println("\t\tExpected output:")
	fmt.Println(expected)
	fmt.Println("\n\t\tActual output:")
	fmt.Println(actual)
	logString += "tester: FAILURE: " + s + "\n" +
		"\t\tExpected output:\n" +
		expected +
		"\n\t\tActual output:\n" +
		actual
}

func logSuccess(s string) {
	fmt.Println("tester: SUCCESS: " + s)
	logString += "tester: SUCCESS: " + s + "\n"
}

func addTabs(s string) string {
	ss := strings.Split(s, "\n")
	for i, s := range ss {
		ss[i] = "\t" + s
	}
	return strings.Join(ss, "\n")
}
