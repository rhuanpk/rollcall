package recorder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rollcall/configs"
	"rollcall/internal/lists"
	"slices"
	"strings"
	"time"
)

const fileName = "record.txt"

var folderName = filepath.Join(configs.FolderAssets, "records")

// File is the global record file to be closed at end of the main.
var File *os.File

// Exec create or open record file.
func Exec() {
	var err error

	if err := os.MkdirAll(folderName, 0775); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalln("error creating records folder:", err)
	}

	File, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("error creating or opening record file:", err)
	}

	if _, err := File.WriteString(time.Now().Format(time.DateTime) + "\n"); err != nil {
		log.Fatalln("error writing record datetime:", err)
	}
}

// Process execute all necessarie final steps of record and must be called at end of the main.
func Process() {
	defer File.Close()
	var names strings.Builder

	var sorted []string
	lists.List.Range(func(key, _ any) bool {
		sorted = append(sorted, key.(string))
		return true
	})
	slices.Sort(sorted)

	for _, name := range sorted {
		value, _ := lists.List.Load(name)
		present := value.(bool)
		fmt.Fprintf(&names, "%-*s: [%s]\n", lists.Max, name, map[bool]string{true: "v", false: "x"}[present])
	}

	if _, err := File.WriteString(names.String()); err != nil {
		log.Println("error writing record names:", err)
	}

	if err := os.Rename(fileName, filepath.Join(
		folderName, time.Now().Format(
			fmt.Sprintf("%s_%s.txt", time.DateOnly, time.TimeOnly),
		),
	)); err != nil {
		log.Println("error changing record file name:", err)
	}

	log.Println("list of presence")
	println(strings.TrimSuffix(names.String(), "\n"))
}
