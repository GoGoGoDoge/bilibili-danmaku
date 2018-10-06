package danmaku

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

type BiliComment struct {
	XMLName    xml.Name  `xml:"i"`
	ChatServer string    `xml:"chatserver"`
	Cid        int64     `xml:"chatid"`
	Mission    string    `xml:"mission"`
	MaxLimit   string    `xml:"maxlimit"`
	State      string    `xml:"state"`
	RealName   string    `xml:"real_name"`
	Source     string    `xml:"source"`
	Comments   []Comment `xml:"d"`
}

type Comment struct {
	Text     string `xml:",chardata"`
	Property string `xml:"p,attr"`
}

type RespComment struct {
	Aid      int64    `json:"aid"`
	Title    string   `json:"title"`
	Cid      int64    `json:"cid"`
	Duration int64    `json:"duration"`
	Danmakus Danmakus `json:"danmaku"`
}

type Danmaku struct {
	stime  int64
	Start  string `json:"start"`
	End    string `json:"end"`
	Effect string `json:"effect"`
	Color  string `json:"color"`
	Text   string `json:"text"`
}

// Danmakus implements sort.Interface.
type Danmakus []Danmaku

func (d Danmakus) Len() int           { return len(d) }
func (d Danmakus) Less(i, j int) bool { return d[i].stime < d[j].stime }
func (d Danmakus) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

type Info struct {
	Stat     InfoStat `json:"stat"`
	Title    string   `json:"title"`
	Duration int64    `json:"duration"`
	Ctime    int64    `json:"ctime"`
}

type InfoStat struct {
	Aid int64 `json:"aid"`
}

func ParseXML(in []byte) (*BiliComment, error) {
	biliComment := &BiliComment{}
	err := xml.Unmarshal(in, biliComment)
	return biliComment, err
}

func ParseInfo(in []byte) (*Info, error) {
	info := &Info{}
	err := json.Unmarshal(in, info)
	return info, err
}

func Convert2RespComment(blComment *BiliComment, info *Info) (*RespComment, error) {
	if blComment == nil || info == nil {
		return nil, fmt.Errorf("cannot convert nil BiliComment or nil Info")
	}
	respComment := &RespComment{
		Aid:      info.Stat.Aid,
		Cid:      blComment.Cid,
		Title:    info.Title,
		Duration: info.Duration,
	}

	// Fill in danmaku format
	danmakus := make(Danmakus, 0)
	for idx, comment := range blComment.Comments {
		props := strings.Split(comment.Property, ",")
		if len(props) != 8 {
			log.Printf("missing property[%d]: %+v\n", idx, props)
			continue
		}
		var danmaku Danmaku
		var err error

		danmaku.stime, danmaku.Start, err = convertTimestamp(props[0], float64(0))
		if err != nil {
			log.Println("time parse error:", err)
			continue
		}
		_, danmaku.End, _ = convertTimestamp(props[0], float64(10))
		danmaku.Text = comment.Text
		danmaku.Color = getColor(props[3])
		danmaku.Effect = getEffect(props[1])
		danmakus = append(danmakus, danmaku)
	}

	sort.Sort(danmakus)
	respComment.Danmakus = danmakus

	return respComment, nil
}

func getEffect(in string) string {
	switch in {
	case "1":
		return "move"
	case "4":
		return "pos_down"
	case "5":
		return "pos_up"
	default:
		return "unknown"
	}
}

func getColor(in string) string {
	bgr, err := strconv.Atoi(in)
	if err != nil {
		return "0xFFFFFF"
	}
	color := ((bgr >> 16) & 255) | (bgr & (255 << 8)) | ((bgr & 255) << 16)

	return fmt.Sprintf("0x%X", color)
}

func convertTimestamp(timestamp string, delta float64) (int64, string, error) {
	f, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		return -1, "", err
	}
	intf := int64(math.Round((f + delta) * 100))
	hour, minute := divmod(intf, 360000)
	minute, second := divmod(minute, 6000)
	second, centsecond := divmod(second, 100)
	return intf, fmt.Sprintf("%d:%02d:%02d.%02d", hour, minute, second, centsecond), nil
}

// length is the text length, ascii char is 0.5 unit where others are 1.
func length(s string) float64 {
	l := 0.0
	for _, r := range s {
		if r < 127 {
			l += 0.5
		} else {
			l += 1
		}
	}
	return l
}

func divmod(f int64, base int64) (a, b int64) { return f / base, f % base }
