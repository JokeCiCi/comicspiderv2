package comic

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"io/ioutil"

	"github.com/JokeCiCi/comicspiderv2/mhttp"
)

func MkDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal("mkdir failed, err:", err)
		}
	}
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ProcessString(old string) (new string) {
	new = strings.ReplaceAll(old, "[?:]", "")
	return
}

func ComicStore(comic *Comic) {
	MkDir(comicRootDir)
	// 创建漫画目录
	comicTitle := ProcessString(comic.Title)
	log.Printf("下载漫画:%s....", comicTitle)

	comicDir := path.Join(comicRootDir, comicTitle)
	MkDir(comicDir)

	comicCoverPath := fmt.Sprintf("%s/%s.jpg", comicRootDir, comicTitle)
	if !Exists(comicCoverPath) {
		mhttp.HttpGetAndStore(comic.CoverURL, comicCoverPath)
	}

	for _, ch := range comic.Chapters {
		// 创建章节目录
		chapterTitle := ProcessString(ch.ChapterTitle)
		chapterDir := path.Join(comicDir, chapterTitle)
		MkDir(chapterDir)

		log.Printf("下载章节:%s....", chapterTitle)
		// 下载章节封面
		chapterCoverPath := fmt.Sprintf("%s/%s.jpg", comicDir, chapterTitle)
		if !Exists(chapterCoverPath) {
			mhttp.HttpGetAndStore(ch.ChapterCoverURL, chapterCoverPath)
		}

		// 下载章节漫画
		i := 1
		for {
			imageURL := fmt.Sprintf("%s%d.jpg", ch.Image.ImagePrefixURL, i)
			imagePath := fmt.Sprintf("%s/%d.jpg", chapterDir, i)
			if !Exists(imagePath) {
				ok := mhttp.HttpGetAndStore(imageURL, imagePath)
				if !ok {
					break
				}
			}
			i++
		}
	}
	log.Printf("下载漫画:%s成功", comicTitle)
}

func InitComics() {
	files, _ := ioutil.ReadDir(comicRootDir)
	for _, f := range files {
		if f.IsDir() {
			comicName := f.Name()
			comicPath := path.Join(comicRootDir, comicName)
			comicURL := path.Join(comicName, "list")
			c, ok := comics[comicName]
			if ok {
				c.ComicPath = comicPath
				c.ComicURL = comicURL
			} else {
				comics[comicName] = &ComicObj{
					ComicName:  comicName,
					ComicPath:  comicPath,
					ComicURL:   comicURL,
					ChapterMap: make(map[string]*ChapterObj),
				}
			}
		} else {
			comicName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			comicCoverURL := path.Join(comicRootDir, f.Name())
			c, ok := comics[comicName]
			if ok {
				c.ComicCoverURL = comicCoverURL
			} else {
				comics[comicName] = &ComicObj{
					ComicCoverURL: comicCoverURL,
					ChapterMap:    make(map[string]*ChapterObj),
				}
			}
		}
	}
	initChapters()
}

func initChapters() {
	for _, c := range comics {
		files, _ := ioutil.ReadDir(c.ComicPath)
		for _, f := range files {
			if f.IsDir() {
				chapterName := f.Name()
				chapterPath := path.Join(c.ComicPath, chapterName)
				chapterURL := path.Join(c.ComicName, chapterName, "list")
				ch, ok := c.ChapterMap[chapterName]
				if ok {
					ch.ChapterPath = chapterPath
					ch.ChapterURL = chapterURL
				} else {
					c.ChapterMap[chapterName] = &ChapterObj{
						ChapterName: chapterName,
						ChapterPath: chapterPath,
						ChapterURL:  chapterURL,
					}
				}
			} else {
				chapterName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
				chapterCoverURL := path.Join(comicRootDir, c.ComicName, f.Name())
				ch, ok := c.ChapterMap[chapterName]
				if ok {
					ch.ChapterCoverURL = chapterCoverURL
				} else {
					c.ChapterMap[chapterName] = &ChapterObj{
						ChapterCoverURL: chapterCoverURL,
					}
				}
			}
		}
	}
	initContents()
}

func initContents() {
	for _, c := range comics {
		for _, ch := range c.ChapterMap {
			files, _ := ioutil.ReadDir(ch.ChapterPath)
			for _, f := range files {
				if !f.IsDir() {
					contentURL := path.Join(ch.ChapterPath, f.Name())
					ch.Contents = append(ch.Contents, contentURL)
				}
			}
		}
	}
}

func ComicObjList() map[string]*ComicObj {
	return comics
}

func ChapterObjList(comicName string) (chs []*ChapterObj) {
	for _, ch := range comics[comicName].ChapterMap {
		chs = append(chs, ch)
	}
	return
}

func ChapterContents(comicName, chapterName string) (chapterContents []string) {
	for _, ch := range comics[comicName].ChapterMap {
		if ch.ChapterName == chapterName {
			return ch.Contents
		}
	}
	return
}
