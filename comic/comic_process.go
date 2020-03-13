package comic

import "strings"

func ProcessComicListPage(html string) (comics []*Comic) {
	comicInfos := comicInfoRe.FindAllStringSubmatch(html, -1)
	for _, v := range comicInfos {
		c := NewComic(v[1], v[2], v[3], v[4])
		comics = append(comics, c)
	}
	return
}

func ProcessComicPage(html string) (chapters []*Chapter) {
	chapterInfos := chapterInfoRe.FindAllStringSubmatch(html, -1)
	for _, v := range chapterInfos {
		chapter := NewChapter(v[1], v[2], v[3])
		chapters = append(chapters, chapter)
	}
	return
}

func ProcessChapterPage(html string) (image *Image) {
	imageURLs := imageURLRe.FindStringSubmatch(html)
	if len(imageURLs) == 0 {
		return
	}
	imagePrefixURL := string([]rune(imageURLs[1])[:strings.LastIndex(imageURLs[1], "/")+1])
	image = NewImage(imagePrefixURL)
	return
}
