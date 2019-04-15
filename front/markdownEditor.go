package main

import (
	md "github.com/shurcooL/github_flavored_markdown"
)

func markdownValue(value []byte) (newValue []byte) {
	return md.Markdown(value)
}
