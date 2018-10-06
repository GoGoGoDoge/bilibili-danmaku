package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bilibili-danmaku/danmaku"
	"github.com/bilibili-danmaku/search"
	"github.com/bilibili-danmaku/util"
	"github.com/gin-gonic/gin"
)

const (
	// FilePrefix is the asset path.
	FilePrefix = "/Users/marco/Documents/Go_Workspace/src/github.com/bilibili-danmaku/util/"
)

// Danmaku handles danmaku request.
func Danmaku(c *gin.Context) {
	title := c.DefaultQuery("title", "")
	incomingUrl := c.DefaultQuery("url", "")

	// Log shall be output to file.
	log.Println("Incoming url is:", incomingUrl)

	errHandler := func() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title cannot be empty",
		})
	}

	if title == "" {
		errHandler()
		return
	}

	simpleTitle, err := util.Trad2SimpleConvert(title)
	if err != nil {
		errHandler()
		return
	}

	log.Println("Simple title:", simpleTitle)

	// matchTitle := search.SearchHelper.GetExactEpisode(simpleTitle)
	matchTitle := search.SearchEngine.GetExactEpisode(simpleTitle)
	// aid := search.SearchHelper.GetAidByTitle(matchTitle)
	aid := search.SearchEngine.GetAidByTitle(matchTitle)
	strAid := strconv.FormatInt(aid, 10)

	log.Println("Mathched title:", matchTitle)
	log.Println("Debug strAid:", strAid)

	possibleRawFiles := []string{
		filepath.Join(FilePrefix, strAid, strAid+"-1.raw"),
		filepath.Join(FilePrefix, "md"+strAid, "md"+strAid+"-1.raw"),
		filepath.Join(FilePrefix, "ss"+strAid, "ss"+strAid+"-1.raw"),
	}

	possibleInfoFiles := []string{
		filepath.Join(FilePrefix, strAid, strAid+"-1.info"),
		filepath.Join(FilePrefix, "md"+strAid, "md"+strAid+"-1.info"),
		filepath.Join(FilePrefix, "ss"+strAid, "ss"+strAid+"-1.info"),
	}

	for idx, rawFile := range possibleRawFiles {
		if _, err := os.Stat(rawFile); !os.IsNotExist(err) {
			xml, err := util.DecompressFile(rawFile)
			if err != nil {
				log.Println("Encounter error", err)
				continue
			}
			bc, err := danmaku.ParseXML(xml)
			if err != nil {
				log.Println("Parse XML error:", err)
				continue
			}
			f, err := os.Open(possibleInfoFiles[idx])
			if err != nil {
				log.Println("Failed to find info:", err)
				continue
			}
			bytes, _ := ioutil.ReadAll(f)
			info, err := danmaku.ParseInfo(bytes)
			if err != nil {
				log.Println("Parse info error:", err)
				continue
			}

			rc, err := danmaku.Convert2RespComment(bc, info)
			if err != nil {
				log.Println("convert 2 resp error:", err)
				continue
			}

			c.JSON(http.StatusOK, rc)
			return
		}
	}
}
