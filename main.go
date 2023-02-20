package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type JdbcConfig struct {
	Address        string
	Port           string
	Database       string
	ConnectionArgs string
}

func main() {
	file, err := os.Open("global.properties")
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := make(map[string]string)
	jdbcConfig := JdbcConfig{}

	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments and empty lines
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Split line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid line:", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Check if key starts with "jdbc.url"
		if strings.HasPrefix(key, "jdbc.url") {
			value = strings.ReplaceAll(value, `\`, "")
			re := regexp.MustCompile(`jdbc:mysql://([^:]+):(\d+)/(.+)`)
			match := re.FindStringSubmatch(value)
			if len(match) < 4 {
				fmt.Println("Invalid jdbc.url:", value)
				continue
			}
			jdbcConfig.Address = match[1]
			jdbcConfig.Port = match[2]
			jdbcConfig.Database = match[3]

			// Remove connection string from database name
			parts := strings.Split(jdbcConfig.Database, "?")
			jdbcConfig.Database = parts[0]

			// Save connection string separately
			if len(parts) > 1 {
				jdbcConfig.ConnectionArgs = "?" + parts[1]
			}
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Print the result
	fmt.Println(config)
	fmt.Println(jdbcConfig)
	fmt.Println(jdbcConfig.Database)
}

//gpt powered
