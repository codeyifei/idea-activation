package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 激活链接
const fileUrl = "http://idea.medeming.com/jets/images/jihuoma.zip"

var NotFoundError = errors.New("未找到激活码")

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	var useOldVersion, isShow bool
	flag.BoolVar(&useOldVersion, "old", false, "是否使用旧版本激活码")
	flag.BoolVar(&isShow, "show", false, "是否显示激活码")
	flag.Parse()

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(os.TempDir())
	buffer.WriteString("activation.zip")
	filename := buffer.String()
	check(download(filename))
	defer func() { _ = os.Remove(filename) }()

	code, err := readCode(filename, useOldVersion)
	check(err)
	check(clipboard.WriteAll(string(code)))
	if useOldVersion {
		color.Green("✓ 旧版激活码已复制到剪切板")
	} else {
		color.Green("✓ 激活码已复制到剪切板")
	}
	if isShow {
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
