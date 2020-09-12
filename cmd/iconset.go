// +build ignore
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	colors = map[string]string{
		"orange": "#ff7305",
		"white":  "#ffffff",
		"blue":   "#14d0f0",
		"red":    "#bb070e",
	}
	offColorIcons = map[string]string{
		"table-v":       "blue",
		"cowboy":        "red",
		"delete-record": "red",
	}
	size = "48"
	tpl  = "\t\"%s\": %s,\n"
	pkg  = `
package assets
var Iconset%s = map[string][]byte{
`
)

func init() {
	rootCmd.AddCommand(iconsetCmd)
}

var iconsetCmd = &cobra.Command{
	Use: "iconset",
	RunE: func(cmd *cobra.Command, args []string) error {
		// export PNG
		err := filepath.Walk("assets/data/iconset", func(path string, info os.FileInfo, err error) error {
			if path == "assets/data/iconset" {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}

			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}

			name := strings.TrimSuffix(info.Name(), ".svg")

			color, ok := offColorIcons[name]
			if !ok {
				color = colors["orange"]
			}

			replacer := strings.NewReplacer(
				"{{STYLE}}",
				fmt.Sprintf(`stroke:none;fill:%s;fill-opacity:1`, color),
				"{{SIZE}}",
				size+"px",
			)

			err = ioutil.WriteFile("assets/iconset-temp/"+info.Name(), []byte(replacer.Replace(string(b))), 0644)
			if err != nil {
				panic(err)
			}

			//code = code + fmt.Sprintf(tpl, strings.TrimSuffix(info.Name(), ".svg"), replacer.Replace(string(b)))
			if strings.HasPrefix(info.Name(), "method-") {
				exportPNG(
					"assets/iconset-temp/"+info.Name(),
					"assets/iconset-temp/"+strings.Replace(info.Name(), ".svg", size+".png", -1),
					"105", "18")
			} else {
				exportPNG(
					"assets/iconset-temp/"+info.Name(),
					"assets/iconset-temp/"+strings.Replace(info.Name(), ".svg", size+".png", -1),
					size, size)
			}

			return nil
		})
		if err != nil {
			panic(err)
		}

		// export to Go
		code := fmt.Sprintf(pkg, size)
		err = filepath.Walk("assets/iconset-temp", func(path string, info os.FileInfo, err error) error {
			if path == "assets/iconset-temp" {
				return nil
			}
			if !strings.HasSuffix(path, size+".png") {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}

			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}

			replacer := strings.NewReplacer(
				"[", "{",
				"]", "}",
				" ", ", ",
			)
			data := replacer.Replace(fmt.Sprintf("%v", b))
			code = code + fmt.Sprintf(tpl, strings.TrimSuffix(info.Name(), size+".png"), data)
			return nil
		})

		code = code + "}"

		err = ioutil.WriteFile("assets/iconset.go", []byte(code), 0644)
		if err != nil {
			return err
		}
		return nil
	},
}

func exportPNG(src, dst, w, h string) {
	goExecutable, err := exec.LookPath("inkscape")
	if err != nil {
		return
	}

	cmd := &exec.Cmd{
		Path: goExecutable,
		//"inkscape -z -w %d -h %d %s.svg -e %s.png"
		Args: []string{goExecutable, "-w", w, "-h",
			h, src, "--export-filename", dst},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	cmd.Process.Release()
}

func exportNoSizePNG(src, dst string) {
	goExecutable, err := exec.LookPath("inkscape")
	if err != nil {
		return
	}

	cmd := &exec.Cmd{
		Path: goExecutable,
		//"inkscape -z -w %d -h %d %s.svg -e %s.png"
		Args:   []string{goExecutable, src, "--export-filename", dst},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	cmd.Process.Release()
}
