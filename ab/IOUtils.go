package ab

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// IOUtils struct for managing file I/O operations.
type IOUtils struct {
	reader *bufio.Reader
	writer *bufio.Writer
	file   *os.File
}

// NewIOUtils creates a new IOUtils instance for reading or writing a file.
func NewIOUtils(filePath, mode string) *IOUtils {
	ioutils := &IOUtils{}
	var err error
	if mode == "read" {
		ioutils.file, err = os.Open(filePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}
		ioutils.reader = bufio.NewReader(ioutils.file)
	} else if mode == "write" {
		ioutils.file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}
		ioutils.writer = bufio.NewWriter(ioutils.file)
	}
	return ioutils
}

// ReadLine reads a line from the file.
func (ioutils *IOUtils) ReadLine() string {
	line, err := ioutils.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ""
	}
	return strings.TrimSpace(line)
}

// WriteLine writes a line to the file.
func (ioutils *IOUtils) WriteLine(line string) {
	_, err := ioutils.writer.WriteString(line + "\n")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	ioutils.writer.Flush()
}

// Close closes the file.
func (ioutils *IOUtils) Close() {
	if ioutils.file != nil {
		ioutils.file.Close()
	}
}

// WriteOutputTextLine prints a prompt and text to the console.
func WriteOutputTextLine(prompt, text string) {
	fmt.Printf("%s: %s\n", prompt, text)
}

// ReadInputTextLine reads a line of input from the console.
func ReadInputTextLine() string {
	return ReadInputTextLineWithPrompt("")
}

// ReadInputTextLineWithPrompt reads a line of input from the console with a prompt.
func ReadInputTextLineWithPrompt(prompt string) string {
	if prompt != "" {
		fmt.Print(prompt + ": ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// ListFiles lists all files in a directory.
func ListFiles(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

// System executes a system command and returns its output.
func Utils_System(evaluatedContents, failedString string) string {
	cmd := exec.Command("sh", "-c", evaluatedContents)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return failedString
	}
	return string(output)
}

// EvalScript evaluates a JavaScript script using the Goja JavaScript engine.
func EvalScript(engineName, script string) (string, error) {
	// Note: For JavaScript evaluation, you'll need to use a library like "github.com/dop251/goja".
	// This placeholder demonstrates the intention.
	return "", fmt.Errorf("JavaScript evaluation not implemented in this example")
}
