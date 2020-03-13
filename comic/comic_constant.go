package comic

import (
	"regexp"
)

const (
	comicRootDir            = "resources/漫画"
	mainEncryptURL          = "https://hhmh109.com/cache/index.html" // 首页加密
	comicEncryptPrefixURL   = "https://hhmh109.com/delistbak.html"
	chapterEncryptPrefixURL = "https://hhmh109.com/detest.php?"
	comicInfoReStr          = `<div\s*class="slmsa"[\s\S]+?href="([\s\S]+?)"[\s\S]+?src="([\s\S]+?)"[\s\S]+?class="t">\s*([\s\S]+?)\s*</span[\s\S]+?<p>\s*([\s\S]+?)\s*</p>`
	chapterInfoReStr        = `<li><a[\s\S]+?href="[\s\S]+?\?([\s\S]+?)"[\s\S]+?"(https[\s\S]+?)"[\s\S]+?pull-left">\s*([\s\S]+?)\s*</?span[\s\S]+?</li>`
	imageURLReStr           = `<img\s*class="lazy"\s*data-original="\s*([\s\S]+?)\s*"[\s\S]+?>`
)

var (
	comicInfoRe   = regexp.MustCompile(comicInfoReStr)
	chapterInfoRe = regexp.MustCompile(chapterInfoReStr)
	imageURLRe    = regexp.MustCompile(imageURLReStr)
	comics        = make(map[string]*ComicObj)
)
