go-yaml-sort
============

# About

`yaml-sort` sorts [YAML](https://yaml.org/) files alphabetically.

(Inspired by [yaml-sort](https://github.com/ddebin/yaml-sort))

# Installation

```shell
git clone git@github.com:ddebin/yaml-sort.git
cd yaml-sort
go build .
```

# Usage

```txt
Usage: yaml-sort [options]
Options:
  -check
    	Check if the given file(s) are already sorted
  -indent int
    	Indentation width to use (in spaces) (default 2)
  -input string
    	The YAML file(s) which need to be sorted (default "-")
  -lineWidth int
    	Wrap line width (-1 for unlimited width) (default 1000)
  -output string
    	The YAML file to output sorted content to
  -quotingStyle string
    	Strings will be quoted using this quoting style (default "single")
  -stdout
    	Output the proposed sort to STDOUT only
Examples:
   yaml-sort --input config.yml
   yaml-sort --input config.yml --lineWidth 100 --stdout
   yaml-sort --input config.yml --indent 4 --output sorted.yml
```
