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

func ComicObjList() map[string]*ComicObj {
	m := make(map[string]*ComicObj)
	files, _ := ioutil.ReadDir(comicRootDir)
	for _, f := range files {
		if f.IsDir() {
			comicName := f.Name()
			comicPath := path.Join(comicRootDir, comicName)
			c, ok := m[comicName]
			if ok {
				c.ComicPath = comicPath
			} else {
				m[comicName] = &ComicObj{
					ComicName: comicName,
					ComicPath: comicPath,
				}
			}
		} else {
			comicName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			comicCoverPath := path.Join(comicRootDir, f.Name())
			c, ok := m[comicName]
			if ok {
				c.ComicCoverPath = comicCoverPath
			} else {
				m[comicName] = &ComicObj{
					ComicCoverPath: comicCoverPath,
				}
			}
		}
	}
	return m
}

func ChapterObjList(comic *ComicObj) map[string]*ChapterObj {
	m := make(map[string]*ChapterObj)
	files, _ := ioutil.ReadDir(comic.ComicPath)
	for _, f := range files {
		if f.IsDir() {
			// 取
			chapterName := f.Name()
			chapterPath := path.Join(comic.ComicPath, chapterName)
			c, ok := m[chapterName]
			if ok {
				c.ChapterPath = chapterPath
			} else {
				m[chapterName] = &ChapterObj{
					ChapterName: chapterName,
					ChapterPath: chapterPath,
				}
			}
		} else {
			chapterName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			chapterCoverPath := path.Join(comic.ComicPath, f.Name())
			c, ok := m[chapterName]
			if ok {
				c.ChapterCoverPath = chapterCoverPath
			} else {
				m[chapterName] = &ChapterObj{
					ChapterCoverPath: chapterCoverPath,
				}
			}
		}
	}
	return m
}

func ChapterContents(chapter *ChapterObj) (chapterContents []string) {
	files, _ := ioutil.ReadDir(chapter.ChapterPath)
	for _, f := range files {
		if !f.IsDir() {
			chapterContent := path.Join(chapter.ChapterPath, f.Name())
			chapterContents = append(chapterContents, chapterContent)
		}
	}
	return
}
