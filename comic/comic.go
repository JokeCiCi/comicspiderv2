package comic

import "fmt"

type Comic struct {
	ComicEncryptUrl string
	CoverURL        string
	Title           string
	Detail          string
	Chapters        []*Chapter
}

type Chapter struct {
	ChapterContentEncryptURL string
	ChapterCoverURL          string
	ChapterTitle             string
	*Image
}

type Image struct {
	ImagePrefixURL string
}

func NewComic(comicEncryptUrl, coverURL, title, detail string) *Comic {
	return &Comic{
		ComicEncryptUrl: comicEncryptPrefixURL + comicEncryptUrl,
		CoverURL:        coverURL,
		Title:           title,
		Detail:          detail,
		Chapters:        make([]*Chapter, 0),
	}
}

func NewChapter(chapterContentEncryptURL, chapterCoverURL, chapterTitle string) *Chapter {
	return &Chapter{
		ChapterContentEncryptURL: chapterEncryptPrefixURL + chapterContentEncryptURL,
		ChapterCoverURL:          chapterCoverURL,
		ChapterTitle:             chapterTitle,
	}
}

func NewImage(imagePrefixURL string) *Image {
	return &Image{
		ImagePrefixURL: imagePrefixURL,
	}
}
func (c *Comic) String() string {
	return fmt.Sprintf("ComicEncryptUrl:[%s] CoverURL:[%s] Title:[%s] Detail:[%s] Chapters:[%s]", c.ComicEncryptUrl, c.CoverURL, c.Title, c.Detail, c.Chapters)
}

func (c *Chapter) String() string {
	return fmt.Sprintf("\n%s %s %s", c.ChapterContentEncryptURL, c.ChapterCoverURL, c.ChapterTitle)
}

func (i *Image) String() string {
	return i.ImagePrefixURL
}

type ComicObj struct {
	ComicName     string
	ComicPath     string
	ComicURL      string
	ComicCoverURL string
	ChapterMap    map[string]*ChapterObj
}

type ChapterObj struct {
	ChapterName     string
	ChapterPath     string
	ChapterURL      string
	ChapterCoverURL string
	Contents        []string
}

func (c *ComicObj) String() string {
	return fmt.Sprintf("ComicName:[%s] ComicPath:[%s] ComicURL:[%s] ComicCoverURL:[%s]", c.ComicName, c.ComicPath, c.ComicURL, c.ComicCoverURL)
}

func (c *ChapterObj) String() string {
	return fmt.Sprintf("ChapterName:[%s] ChapterPath:[%s] ChapterURL:[%s] ChapterCoverURL:[%s] Contents:[%s]", c.ChapterName, c.ChapterPath, c.ChapterURL, c.ChapterCoverURL, c.Contents)
}
