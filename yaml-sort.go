package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	var (
		input        = flag.String("input", "-", "The YAML file(s) which need to be sorted")
		output       = flag.String("output", "", "The YAML file to output sorted content to")
		stdout       = flag.Bool("stdout", false, "Output the proposed sort to STDOUT only")
		check        = flag.Bool("check", false, "Check if the given file(s) are already sorted")
		indent       = flag.Int("indent", 2, "Indentation width to use (in spaces)")
		quotingStyle = flag.String("quotingStyle", "single", "Strings will be quoted using this quoting style")
		lineWidth    = flag.Int("lineWidth", 1000, "Wrap line width (-1 for unlimited width)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "--input config.yml")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "--input config.yml --lineWidth 100 --stdout")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "--input config.yml --indent 4 --output sorted.yml")
		fmt.Fprintln(os.Stderr, "  cat config.yml |", os.Args[0])
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	success := true

	files := flag.Args()
	if *input != "-" {
		files = append(files, *input)
	}

	for _, file := range files {
		isStdin := file == "-"
		var content []byte
		var err error

		if isStdin {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				content, err = ioutil.ReadAll(os.Stdin)
			} else {
				fmt.Fprintln(os.Stderr, "Missing filename or input (\"yaml-sort --help\" for help)")
				os.Exit(22)
			}
		} else {
			content, err = ioutil.ReadFile(file)
		}

		if err != nil {
			success = false
			fmt.Println(err)
			continue
		}

		var outContent []byte

		// Parse YAML
		var data interface{}
		if err := yaml.Unmarshal(content, &data); err != nil {
			success = false
			fmt.Println(err)
			continue
		}

		// Sort keys
		if !*check {
			data = sortKeys(data)
		}

		// Dump YAML
		outContent, err = dumpYAML(data, *indent, *quotingStyle, *lineWidth)
		if err != nil {
			success = false
			fmt.Println(err)
			continue
		}

		// Check for sorting
		if *check {
			if string(outContent) != string(content) {
				success = false
				fmt.Printf("'%s' is not sorted and/or formatted (indent, line width).\n", file)
			}
		} else if *stdout || isStdin && *output == "" {
			fmt.Print(string(outContent))
		} else {
			outputFile := *output
			if outputFile == "" {
				outputFile = file
			}
			if err := ioutil.WriteFile(outputFile, outContent, os.ModePerm); err != nil {
				success = false
				fmt.Println(err)
			}
		}
	}

	if !success {
		os.Exit(1)
	}
}

func dumpYAML(data interface{}, indent int, quotingStyle string, lineWidth int) ([]byte, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	if quotingStyle == "double" {
		yamlData = bytesReplace(yamlData, "'", "\"")
	}

	if lineWidth > 0 {
		lines := strings.Split(string(yamlData), "\n")
		for i, line := range lines {
			if len(line) > lineWidth {
				lines[i] = wrapLine(line, lineWidth, indent)
			}
		}
		return []byte(strings.Join(lines, "\n")), nil
	}

	return yamlData, nil
}

func sortKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		sorted := make(map[interface{}]interface{})
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key.(string))
		}
		sort.Strings(keys)
		for _, key := range keys {
			sorted[key] = sortKeys(v[key])
		}
		return sorted
	case []interface{}:
		for i, item := range v {
			v[i] = sortKeys(item)
		}
	}
	return data
}

func bytesReplace(input []byte, old, new string) []byte {
	return []byte(strings.Replace(string(input), old, new, -1))
}

func wrapLine(line string, lineWidth, indent int) string {
	indentStr := strings.Repeat(" ", indent)
	parts := splitLine(line, lineWidth-len(indentStr))
	for i, part := range parts {
		if i == 0 {
			parts[i] = fmt.Sprintf("%s%s", indentStr, part)
		} else {
			parts[i] = fmt.Sprintf("%s%s%s", indentStr, strings.Repeat(" ", lineWidth), part)
		}
	}
	return strings.Join(parts, "\n")
}

func splitLine(s string, width int) []string {
	var parts []string
	for len(s) > width {
		i := width
		for ; i >= 0; i-- {
			if s[i] == ' ' {
				parts = append(parts, s[:i])
				s = s[i+1:]
				break
			}
		}
		if i == -1 {
			parts = append(parts, s[:width])
			s = s[width:]
		}
	}
	if len(s) > 0 {
		parts = append(parts, s)
	}
	return parts
}
