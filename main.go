package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/JokeCiCi/comicspiderv2/comic"
	"github.com/JokeCiCi/comicspiderv2/mhttp"
)

func WorkJob(jobChan <-chan *comic.Comic, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobChan:
			if !ok {
				break
			}
			comic.ComicStore(job)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func GenerateJob(jobChan chan<- *comic.Comic, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(jobChan)
	i := 1
	for {
		url := fmt.Sprintf("https://hhmh109.com/detoplist.html?order=lastupdate&page=%d", i)
		comicListHTML := mhttp.HTTPEncryptGet(url)
		comics := comic.ProcessComicListPage(comicListHTML)
		if len(comics) == 0 {
			break
		}
		fmt.Printf("url:%s comics:%d \n", url, len(comics))
		for _, c := range comics {
			// fmt.Println("comic url", c.ComicEncryptUrl)
			comicHTML := mhttp.HTTPEncryptGet(c.ComicEncryptUrl)
			chapters := comic.ProcessComicPage(comicHTML)
			for _, ch := range chapters {
				// fmt.Println("chapter url", ch.ChapterContentEncryptURL)
				chapterHTML := mhttp.HTTPEncryptGet(ch.ChapterContentEncryptURL)
				image := comic.ProcessChapterPage(chapterHTML)
				for nil == image {
					chapterHTML = mhttp.HTTPEncryptGet(ch.ChapterContentEncryptURL)
					image = comic.ProcessChapterPage(chapterHTML)
					// log.Println("image is nil", image)
				}
				ch.Image = image
			}
			c.Chapters = chapters
			jobChan <- c
			// PrintComic(c)
		}
		i++
	}
}

func PrintComic(c *comic.Comic) {
	f, _ := os.OpenFile("./info", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	for _, ch := range c.Chapters {
		fmt.Fprintln(f, c.Title, ch.ChapterTitle, ch.ImagePrefixURL)
	}
}

func DownloadComic() {
	var wg sync.WaitGroup
	jobChan := make(chan *comic.Comic, 16)
	wg.Add(1)
	go GenerateJob(jobChan, &wg)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go WorkJob(jobChan, &wg)
	}
	wg.Wait()
}

func StartServe() {

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("resources"))))
	tmpls, err := template.ParseGlob("resources/tmpl/*")
	if err != nil {
		log.Println("template ParseGlob err:", err)
	}
	// // 列出所有漫画
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		m := comic.ComicObjList()
		tmpls.ExecuteTemplate(w, "comic_list.tmpl", m)
	})
	// // 列出所有章节
	// http.HandleFunc("/:chapter/list", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpls.ExecuteTemplate(w, "list.tmpl", nil)
	// 	// fmt.Fprintln(w,"hello")
	// })

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "hello world")
	// })

	http.ListenAndServe(":80", nil)
}

func main() {
	// DownloadComic()
	StartServe()
	// cmap := comic.ComicObjList()
	// for _, c := range cmap {
	// 	// fmt.Println(c)
	// 	chmap := comic.ChapterObjList(c)
	// 	for _, ch := range chmap {
	// 		cns := comic.ChapterContents(ch)
	// 		for _, cn := range cns {
	// 			fmt.Println(cn)
	// 		}
	// 	}
	// }
}
