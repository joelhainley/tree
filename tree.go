package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


const ITEM_GLYPH = "\u251C\u2500"
const LAST_ITEM_GLYPH = "\u2514\u2500"
const LINE_CONT_GLYPH = "\u2502"


// TODO support the ability to pass in separator (default to '/')
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lineCount := 0
	root := Entry{}
	root.name = "/"
	for scanner.Scan() {
		//fmt.Printf("%d: %s\n", lineCount, scanner.Text())
		processEntry(scanner.Text(), "/", &root)
		lineCount++
	}
	var entryPrefix[] string
	printEntry(&root, 0, entryPrefix, false)
}

type Entry struct {
	name string
	children[] *Entry
}

func processEntry(entry string, separator string, root *Entry) {
	parts := strings.Split(entry, separator)
	_, remaining := parts[0], parts[1:]
	processParts(remaining, root)
}

func processParts(parts[] string, currentNode *Entry){
	if len(parts) > 1 {
		part, remaining := parts[0], parts[1:]

		// does the current nodes children contain an entry for part yet?
		found, child := findChildWithName(currentNode.children, part)

		if found {
			// if it does then get that entry and descend to the next level
			processParts(remaining, child)
		} else {
			newEntry := Entry{}
			newEntry.name = part

			currentNode.children = append(currentNode.children, &newEntry)

			processParts(remaining, &newEntry)
		}
	} else if len(parts) == 1 {
		newEntry := Entry{}
		newEntry.name = parts[0]
		currentNode.children = append(currentNode.children, &newEntry)
	}
}

func findChildWithName(children[] *Entry, name string) (bool, *Entry) {
	for _, child := range children {
		if child.name == name {
			return true, child
		}
	}
	return false, &Entry{}
}

func dumpEntries(entry *Entry, depth int){
	fmt.Printf("depth: %d > child name: %s\n", depth, entry.name)
	for _, child := range entry.children {
		dumpEntries(child, depth+1)
	}
}

func printEntry(entry *Entry, curDepth int, prefix[] string, isLastEntry bool){
	glyph := ITEM_GLYPH

	if isLastEntry {
		glyph = LAST_ITEM_GLYPH
	}

	if curDepth == 0 {
		fmt.Printf("%s\n", entry.name)
	} else {
		fmt.Printf("%s%s %s\n", getPrefixSlug(prefix), glyph, entry.name)
	}

	for index, child := range entry.children {
		lastEntry := false
		if index + 1 == len(entry.children){
			lastEntry = true
		}

		var nextPrefix[] string

		if isLastEntry && lastEntry {
			nextPrefix = append(prefix, "")
		} else {
			if isLastEntry {
				nextPrefix = append(prefix, "" /*LINE_CONT_GLYPH + "b" */)
			} else {
				nextPrefix = append(prefix, LINE_CONT_GLYPH)
			}
		}

		if curDepth == 0 {
			nextPrefix = prefix
		}

		if len(child.children) > 0 {
			printEntry(child, curDepth+1, nextPrefix , lastEntry)
		} else {
			printEntry(child, curDepth+1, nextPrefix, lastEntry)
		}
	}
}

func getPrefixSlug(prefix[] string) string {
	var sb strings.Builder

	for i := 0; i < len(prefix); i++ {
		sb.WriteString(prefix[i])
		sb.WriteString("   ")
	}
	return sb.String()
}