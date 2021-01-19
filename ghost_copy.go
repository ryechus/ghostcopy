package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usageDocs() {
	fmt.Printf("usage: %s source destination\n", os.Args[0])
}

func main() {
	working_dir := flag.String("w", "", "working directory to execute the copy from")
	flag.Usage = usageDocs
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
	if *working_dir != "" {
		err := os.Chdir(*working_dir)
		if err != nil {
			fmt.Printf("%s does not exist \n", *working_dir)
		}
		cwd, _ = os.Getwd()
		fmt.Printf("working dir is %s\n", cwd)
	}
	source := flag.Arg(0)
	destination := flag.Arg(1)
	source_dir := filepath.Join(cwd, source)

	err := filepath.Walk(source_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		rep_path := strings.Replace(path, source_dir, "", -1)
		new_path := filepath.Join(cwd, "../", destination, source, rep_path)
		if path == new_path {
			fmt.Printf("Can't copy because the paths are the same: %s\n", path)
			return nil
		}

		if info.IsDir() {
			os.MkdirAll(new_path, os.ModePerm)
		} else {
			f, _ := os.Create(new_path)
			f.Close()
		}

		fmt.Printf("Created ghost copy of %s at %s\n", path, new_path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}
}
