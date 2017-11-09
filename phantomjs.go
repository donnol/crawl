package main

import (
	"os/exec"
)

// 调用命令 phantomjs [--script-encoding=encoding --output-encoding=gbk] request.js {url}
// 如果出现中文无显示，请安装字体 sudo apt-get install xfonts-wqy
func phantom(url string) ([]byte, error) {
	cmd := exec.Command("phantomjs", "request.js", url)
	output, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}
	return output, nil
}
