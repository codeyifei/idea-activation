package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 激活链接
// const fileUrl = "http://idea.medeming.com/jets/images/jihuoma.zip"
const fileUrl = "http://idea.medeming.com/a/jihuoma1.zip"

var NotFoundError = errors.New("未找到激活码")

func main() {
	defer pause()
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	var useOldVersion, isShow, url, noCopy bool
	flag.BoolVar(&useOldVersion, "old", false, "是否使用旧版本激活码")
	flag.BoolVar(&isShow, "show", false, "是否显示激活码")
	flag.BoolVar(&url, "url", false, "是否显示激活链接")
	flag.BoolVar(&noCopy, "no-copy", false, "不复制到剪切板")
	flag.Parse()

	if url {
		color.Green("激活文件链接为：%s", fileUrl)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(os.TempDir())
	buffer.WriteString("activation.zip")
	filename := buffer.String()
	check(download(filename))
	defer func() { _ = os.Remove(filename) }()

	code, err := readCode(filename, useOldVersion)
	check(err)
	if !noCopy {
		check(clipboard.WriteAll(string(code)))
		if useOldVersion {
			color.Green("✓ 旧版激活码已复制到剪切板")
		} else {
			color.Green("✓ 激活码已复制到剪切板")
		}
	} else {
		if useOldVersion {
			color.Green("✓ 旧版激活码已生成")
		} else {
			color.Green("✓ 激活码已生成")
		}
	}
	if isShow || noCopy {
		fmt.Println(string(code))
	}
}

func download(filename string) error {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, resp.Body)
	return err
}

func readCode(filename string, useOldVersion bool) ([]byte, error) {
	rc, err := zip.OpenReader(filename)
	check(err)
	defer func() { _ = rc.Close() }()
	var subs []string
	if useOldVersion {
		subs = []string{"or earlier", "以前", "之前"}
	} else {
		subs = []string{"or later", "以后", "之后"}
	}
	for _, file := range rc.File {
		name, err := toUTF8([]byte(file.Name))
		if err != nil {
			return nil, err
		}
		if isStrContain(string(name), subs...) {
			f, err := file.Open()
			if err != nil {
				return nil, err
			}
			l := bufio.NewReader(f)
			content := make([]byte, 0)
			for {
				c, _, err := l.ReadLine()
				if err == io.EOF {
					break
				}
				content = append(content, c...)
			}
			_ = f.Close()
			return content, nil
		}
	}
	return nil, NotFoundError
}

func isStrContain(target string, subs ...string) bool {
	for _, v := range subs {
		if n := strings.Index(target, v); n > -1 {
			return true
		}
	}
	return false
}

func toUTF8(s []byte) ([]byte, error) {
	if isGBK(s) {
		reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
		return ioutil.ReadAll(reader)
	}
	return s, nil
}

func isGBK(s []byte) bool {
	l := len(s)
	var i = 0
	for i < l {
		if s[i] <= 0x7f {
			i++
			continue
		} else {
			if s[i] >= 0x81 && s[i] <= 0xfe && s[i+1] >= 0x40 && s[i+1] <= 0xfe && s[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func pause() {
	if runtime.GOOS == "windows" {
		fmt.Print("Press Enter or Ctrl-C to exit...")
		bufio.NewScanner(os.Stdin).Scan()
	}
}
