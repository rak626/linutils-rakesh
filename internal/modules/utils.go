package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func appendToFileIfMissing(filePath, line string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create file if it doesn't exist
		return os.WriteFile(filePath, []byte(line+"\n"), 0644)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), line) {
			fmt.Printf("Line already exists in %s: %s\n", filePath, line)
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Append line
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + line + "\n"); err != nil {
		return err
	}

	fmt.Printf("Added to %s: %s\n", filePath, line)
	return nil
}

func prependToFileIfMissing(filePath, line string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create file if it doesn't exist
		return os.WriteFile(filePath, []byte(line+"\n"), 0644)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), line) {
		fmt.Printf("Line already exists in %s: %s\n", filePath, line)
		return nil
	}

	newContent := line + "\n" + string(content)
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Prepended to %s: %s\n", filePath, line)
	return nil
}

func copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, input, 0644)
}
