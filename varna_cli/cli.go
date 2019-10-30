package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "Varna"
	app.Usage = "A CLI for compiling LLVM program to WebAssembly"
	app.Author = "Achala"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name: "Rust",
			// Aliases: []string{"t"},
			Usage: "Rust program",
			Subcommands: []cli.Command{
				{
					Name:  "init",
					Usage: "Create project file tree",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "message, m"},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("Success!!! You created file tree ", c.Args().First())
						directory("contracts/src")
						directory1("contracts")
						return nil
					},
				},
				{
					Name: "compile",
					// Aliases: []string{"rs"},
					Usage: "Compile Rust program",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "message, m"},
					},
					Action: func(c *cli.Context) {
						fmt.Println(c.Generic("m"))
						fmt.Println("Varna API is up and running \n ")

						f, _ := os.Open("/home/achala/go/src/varna/varna_cli/contracts/src/main.rs")
						scanner := bufio.NewReader(f)
						var done bool
						var abc string
						for !done {
							dataString, err := scanner.ReadString('\n')
							if err != nil {
								done = true
							}
							abc = abc + dataString
						}
						test1 := strings.ReplaceAll(abc, "\n", "")
						fmt.Println(test1)
						newfile, err := os.Create("data.rs")
						data, err := io.WriteString(newfile, test1)
						log.Println("Size of the file in bytes:", data)

						data1, err := ioutil.ReadFile("data.rs")
						if err != nil {
							panic(err)
						}

						body := `{"code": "` + string(data1) + `"}`
						resp, err := http.Post("http://localhost:8083/build/rs", "application/json", strings.NewReader(body))
						if err != nil {
							panic(err)
						}
						respBt, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							panic(err)

						}
						fmt.Println("Your Compiled WASM code in json format \n")
						// result := `"output":` + "\n" + string(respBt) + ` `
						wasm := string(respBt)
						// tst := json.Unmarshal(respBt, )

						outfile, err := os.Create("result.wasm")
						file2, err := io.WriteString(outfile, wasm)
						log.Println("size of the output file in bytes:", file2)
					},
				},
			},
		},
		{
			Name: "C",
			// Aliases: []string{"t"},
			Usage: "C program",
			Subcommands: []cli.Command{
				{
					Name:  "init",
					Usage: "Create project file tree",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "message, m"},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("Success!!! You created file tree ", c.Args().First())
						directory2("contracts")
						return nil
					},
				},
				{
					Name: "compile",
					// Aliases: []string{"rs"},
					Usage: "Compile C program",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "message, m"},
					},
					Action: func(c *cli.Context) {
						fmt.Println(c.Generic("m"))
						fmt.Println("Varna API is up and running \n")

						// codedata, err := os.Open("/home/achala/go/src/varna/varna_cli/contracts/src/main.rs")
						// if err != nil {
						// 	panic(err)
						// }
						// defer codedata.Close()

						// newfile, err := os.Create("data.txt")
						// file, err := io.Copy(newfile, codedata)

						f, _ := os.Open("/home/achala/go/src/varna/varna_cli/contracts/main.c")
						scanner := bufio.NewReader(f)
						var done bool
						var abc string
						for !done {

							dataString, err := scanner.ReadString('\n')
							if err != nil {
								done = true
							}
							abc = abc + dataString
						}
						test1 := strings.ReplaceAll(abc, "\n", "")

						newfile, err := os.Create("data.c")
						data, err := io.WriteString(newfile, test1)
						fmt.Println(data)

						data1, err := ioutil.ReadFile("data.c")
						if err != nil {
							panic(err)
						}

						// body := `{"code": "` + c.String("message") + `"}`
						body := `{"code": "` + string(data1) + `"}`
						resp, err := http.Post("http://localhost:8083/build/c", "application/json", strings.NewReader(body))
						if err != nil {
							panic(err)
						}
						respBt, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							panic(err)

						}
						// res1 := `"success":"true"` + "\n"
						fmt.Println("Your Compiled WASM code in json format \n")
						result := `"output":` + "\n" + string(respBt) + ` `
						fmt.Println(string(result))
						outfile, err := os.Create("result.wasm")
						file2, err := io.WriteString(outfile, result)
						log.Println(file2)

					},
				},
			},
		},
	}
}

func directory(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)

	emptyFile, err := os.Create(filepath.Join(dirName, filepath.Base("main.rs")))
	if err != nil {
		log.Fatal(err)
	}
	text := []string{"#[no_mangle]", "pub extern fn sum(x:u64, y:u64) -> u64 {", "   x + y", "}", "fn main(){", "}"}

	for _, data := range text {
		fmt.Fprintln(emptyFile, data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	log.Println(emptyFile)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func directory1(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)

	cargotoml, err := os.Create(filepath.Join(dirName, filepath.Base("Cargo.toml")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cargotoml)
	text := []string{"[package]", "", "name = \"varna\"", "version = \"0.0.1\" ", "authors = [\"Achala Bhatt<achalasbhatt@gmail.com>\"]", "", "[dependencies]"}

	for _, data := range text {
		fmt.Fprintln(cargotoml, data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	cargotoml.Close()

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func directory2(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)

	emptyFile, err := os.Create(filepath.Join(dirName, filepath.Base("main.c")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(emptyFile)
	text := []string{"#include <stdio.h>", "int main() {", "	printf(\"Hello from C!\");", "		return 0;", "}"}

	for _, data := range text {
		fmt.Fprintln(emptyFile, data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	emptyFile.Close()

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
