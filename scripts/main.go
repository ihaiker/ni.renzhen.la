package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type (
	Nav struct {
		Name       string
		Link, Path string
		Order      int
		Step       int
		Items      []*Nav
	}
)

func write(step int, base string, navs []*Nav) *bytes.Buffer {
	lines := bytes.NewBufferString("")
	for _, nav := range navs {
		if nav.Step > 0 && filepath.Base(nav.Link) == "index.md" {
			continue
		}
		prefix := strings.Repeat("    ", nav.Step-step)
		rel, _ := filepath.Rel(base, nav.Path)
		lines.WriteString(fmt.Sprintf("%s* [%s](%s)\n", prefix, nav.Name, rel))
		subLines := write(step, base, nav.Items).String()
		lines.WriteString(subLines)

		if len(nav.Items) > 0 {
			path := filepath.Join(nav.Path, "index.md")
			index := write(nav.Step+1, nav.Path, nav.Items)
			f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0655)
			_, _ = f.WriteString(index.String())
			_ = f.Close()
		}
	}
	return lines
}

func main() {
	dir, _ := os.Getwd()
	docs := filepath.Join(dir, "docs")
	navs := walk(0, docs, docs)

	f, _ := os.OpenFile(filepath.Join(dir, "mkdocs.yml"),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0655)
	f.WriteString("nav: \n")
	for _, nav := range navs {
		f.WriteString(fmt.Sprintf("  - '%s': %s \n", nav.Name, nav.Link))
	}
	_ = f.Close()
	fmt.Println(write(0, docs, navs).String())
}

func walk(step int, base, path string) (navs []*Nav) {
	rel, _ := filepath.Rel(base, path)
	files, _ := ioutil.ReadDir(path)

	for _, file := range files {
		if file.IsDir() {
			metadata := getMetadata(filepath.Join(path, file.Name()), "index.md")
			nav := &Nav{
				Name: metadata["title"], Link: filepath.Join(rel, file.Name()),
				Step: step, Order: 1000, Path: filepath.Join(path, file.Name()),
			}
			nav.Items = walk(step+1, base, filepath.Join(path, file.Name()))
			navs = append(navs, nav)
		} else {
			if filepath.Ext(file.Name()) != ".md" {
				continue
			}
			metadata := getMetadata(path, file.Name())
			nav := &Nav{
				Name: metadata["title"], Link: filepath.Join(rel, file.Name()),
				Step: step, Order: 1000, Path: filepath.Join(path, file.Name()),
			}
			nav.Order, _ = strconv.Atoi(metadata["order"])
			navs = append(navs, nav)
		}
	}
	if len(navs) > 0 {
		sort.Slice(navs, func(i, j int) bool {
			return navs[i].Order < navs[j].Order
		})
	}
	return navs
}

func getMetadata(path, name string) map[string]string {
	file := filepath.Join(path, name)
	metadata := make(map[string]string)

	lines := newIterator(file)
	if _, line, num, has := lines.FindLine("---", "# "); has {
		if line == "---" && num == 0 {
			if args, _, _, has := lines.FindLine("---"); has {
				for _, arg := range args {
					keyVal := strings.SplitN(arg, ": ", 2)
					if unquoteValue, err := strconv.Unquote(keyVal[1]); err == nil {
						metadata[keyVal[0]] = unquoteValue
					} else {
						metadata[keyVal[0]] = keyVal[1]
					}
				}
			}
		} else if strings.HasPrefix(line, "# ") {
			metadata["title"] = line[2:]
		}
		return metadata
	}

	if _, has := metadata["title"]; !has {
		if name == "index.md" {
			metadata["title"] = filepath.Base(path)
		} else {
			metadata["title"] = strings.Replace(name, ".md", "", 1)
		}
	}
	return metadata
}

type lineIterator struct {
	reader  *bufio.Reader
	file    *os.File
	current []byte
}

func (self *lineIterator) FindLine(args ...string) (before []string, line string, lineNumber int, find bool) {
	lineNumber = 0
	for self.HasNext() {
		newLine := string(self.Next())
		for _, arg := range args {
			if strings.HasPrefix(newLine, arg) {
				line = newLine
				find = true
				return
			}
		}
		before = append(before, newLine)
		lineNumber += 1
	}
	return
}
func (self *lineIterator) HasNext() bool {
	line, _, err := self.reader.ReadLine()
	if err != nil {
		self.current = nil
		return false
	}
	self.current = line
	return true
}
func (self *lineIterator) Next() []byte {
	defer func() {
		self.current = nil
	}()
	return self.current
}
func (self *lineIterator) Close() error {
	return self.file.Close()
}
func newIterator(path string) *lineIterator {
	file, _ := os.Open(path)
	return &lineIterator{
		file:   file,
		reader: bufio.NewReader(file),
	}
}
