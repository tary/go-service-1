package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/giant-tech/go-service/base/gameconfig"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("jsonDir", "", "json文件夹，不能为空")
	pflag.String("codeDir", "", "代码生成的文件夹，不能为空")
}

func main() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	jsonDir := viper.GetString("jsonDir")
	fmt.Printf("jsonDir is %s\n", jsonDir)
	if jsonDir == "" {
		fmt.Printf("jsonDir is empty ")
		fmt.Println("Usage:")
		pflag.PrintDefaults()
		return
	}

	codeDir := viper.GetString("codeDir")
	if codeDir == "" {
		fmt.Printf("codeDir is empty ")
		fmt.Println("Usage:")
		pflag.PrintDefaults()
		return
	}

	//先读取alias.json
	aliaspath := jsonDir + "alias.json"
	err := filepath.Walk(aliaspath, func(path string, file os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".json") {
			fmt.Println("path: ", path)

			defInfo := gameconfig.New(path)
			name := defInfo.Get("name")
			var typet *TypeTemplate

			if file.Name() == "alias.json" {
				typet = NewTypeTemplate(name.(string))
				types := defInfo.Get("types").(map[string]interface{})
				for index, info := range types {
					infoMap := info.(map[string]interface{})
					//typ := infoMap["type"]
					typet.AddType(index, infoMap)
				}
			}

			var targetFileName string
			if file.Name() == "alias.json" {
				targetFileName = fmt.Sprintf("%s/alias.go", codeDir)
			}

			fmt.Println("targetFileName: ", targetFileName)

			f, err := os.Create(targetFileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if file.Name() == "alias.json" {
				fmt.Println("file name= ", file.Name())
				_, err = f.WriteString(typet.genString())
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			f.Sync()
			f.Close()

			fmt.Println("Process Done")
		}

		return nil
	})

	err = filepath.Walk(jsonDir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			fmt.Println("path: ", path)

			defInfo := gameconfig.New(path)
			name := defInfo.Get("name")

			var t *Template

			if file.Name() != "alias.json" {
				t = NewTemplate(name.(string))
				props := defInfo.Get("props").(map[string]interface{})
				for prop, info := range props {
					infoMap := info.(map[string]interface{})
					typ := infoMap["type"].(string)
					t.AddType(prop, typ)
				}

				var targetFileName string

				targetFileName = fmt.Sprintf("%s/%sDef.go", codeDir, name)

				fmt.Println("targetFileName: ", targetFileName)

				f, err := os.Create(targetFileName)
				if err != nil {
					fmt.Println(err)
					return err
				}

				_, err = f.WriteString(t.genString())
				if err != nil {
					fmt.Println(err)
					return err
				}
				f.Sync()
				f.Close()
				fmt.Println("Process Done")
			}

		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	/*err = filepath.Walk(jsonDir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			fmt.Println("path: ", path)

			defInfo := gameconfig.New(path)
			name := defInfo.Get("name")

			var t *Template
			var typet *TypeTemplate

			if file.Name() == "alias.json" {
				typet = NewTypeTemplate(name.(string))
				types := defInfo.Get("types").(map[string]interface{})
				for index, info := range types {
					infoMap := info.(map[string]interface{})
					//typ := infoMap["type"]
					typet.AddType(index, infoMap)
				}
			} else {
				t = NewTemplate(name.(string))
				props := defInfo.Get("props").(map[string]interface{})
				for prop, info := range props {
					infoMap := info.(map[string]interface{})
					typ := infoMap["type"].(string)
					t.AddType(prop, typ)
				}
			}

			var targetFileName string
			if file.Name() == "alias.json" {
				targetFileName = fmt.Sprintf("%s/alias.go", codeDir)
			} else {
				targetFileName = fmt.Sprintf("%s/%sDef.go", codeDir, name)
			}
			fmt.Println("targetFileName: ", targetFileName)

			f, err := os.Create(targetFileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if file.Name() == "alias.json" {
				fmt.Println("file name= ", file.Name())
				_, err = f.WriteString(typet.genString())
			} else {
				_, err = f.WriteString(t.genString())
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			f.Sync()
			f.Close()

			fmt.Println("Process Done")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}*/
}
