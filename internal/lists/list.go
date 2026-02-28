package lists

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rollcall/configs"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-hl/normalize"
)

var folderLists = filepath.Join(configs.FolderAssets, "lists")

const Ext = ".txt"

var (
	// List is the global list of presence loaded from file.
	List sync.Map

	// String is the global list of presence concatenaded by new line.
	String string

	// Max is the largest name in the list of presence.
	Max int

	// File is the global name that represent the original list name.
	Name string
)

// Exec load the list of presence from file.
func Exec() {
	var option int

	entries, err := os.ReadDir(folderLists)
	if err != nil {
		log.Fatalln("error listing lists folder:", err)
	}

	if len(entries) <= 0 {
		log.Fatalln("error not found lists files")
	}

	log.Println("lists of presence")
	for index, entry := range slices.Clone(entries) {
		if !strings.HasSuffix(entry.Name(), Ext) {
			entries = slices.Delete(entries, index, index+1)
			continue
		}
		log.Printf("%d. %s", index, entry.Name())
	}

	fmt.Printf("%s option: ", time.Now().Format("2006/01/02 "+time.TimeOnly))
	if _, err := fmt.Scanln(&option); err != nil {
		log.Fatalln("error choosing list file:", err)
	}

	fileName := entries[option].Name()
	filePath := filepath.Join(folderLists, fileName)
	Name = fileName

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("error reading list file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := normalize.String(scanner.Text())

		List.Store(name, false)
		String += name + "\n"

		if max := len(name); max > Max {
			Max = max
		}
	}
}
