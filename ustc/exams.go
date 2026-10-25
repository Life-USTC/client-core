package ustc

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const ExamArrangementURL = "https://jw.ustc.edu.cn/for-std/exam-arrange"

// Exam retains the separate fields used by the two existing table formats.
// DateTime stays textual; native callers retain their date/timezone policy.
type Exam struct {
	CourseName   string `json:"courseName"`
	LessonCode   string `json:"lessonCode"`
	ExamType     string `json:"examType"`
	DateTime     string `json:"dateTime"`
	Location     string `json:"location"`
	Seat         string `json:"seat"`
	Status       string `json:"status"`
	RoomName     string `json:"roomName"`
	BuildingName string `json:"buildingName"`
	District     string `json:"district"`
	Detail       string `json:"detail"`
}

// ParseExamStatusTable retains the CLI table selector, descendant text, column
// mapping and sorting. It does not infer the separate room-column format.
func ParseExamStatusTable(data []byte) ([]Exam, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse exams page: %w", err)
	}
	items := []Exam{}
	doc.Find("table tbody tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 4 {
			return
		}
		values := make([]string, 7)
		for i := 0; i < cells.Length() && i < len(values); i++ {
			values[i] = strings.Join(strings.Fields(cells.Eq(i).Text()), " ")
		}
		if values[0] != "" {
			items = append(items, Exam{CourseName: values[0], LessonCode: values[1], ExamType: values[2], DateTime: values[3], Location: values[4], Seat: values[5], Status: values[6]})
		}
	})
	sort.Slice(items, func(i, j int) bool {
		if items[i].DateTime != items[j].DateTime {
			return items[i].DateTime < items[j].DateTime
		}
		if items[i].CourseName != items[j].CourseName {
			return items[i].CourseName < items[j].CourseName
		}
		return items[i].LessonCode < items[j].LessonCode
	})
	return items, nil
}

// ParseExamRoomTable retains the iOS service's eight-column layout, direct
// cell text and source order. Room, building, district and detail stay separate.
func ParseExamRoomTable(data []byte) ([]Exam, error) {
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("cannot decode exam HTML as UTF-8")
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	items := []Exam{}
	doc.Find("#exams > tbody > tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Children()
		if cells.Length() < 8 {
			return
		}
		values := make([]string, 8)
		for i := range values {
			values[i] = examOwnText(cells.Eq(i).Nodes[0])
		}
		items = append(items, Exam{LessonCode: values[0], ExamType: values[1], CourseName: values[2], DateTime: values[3], RoomName: values[4], BuildingName: values[5], District: values[6], Detail: values[7]})
	})
	return items, nil
}

func examOwnText(cell *html.Node) string {
	var text strings.Builder
	for child := cell.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			text.WriteString(child.Data)
		} else if child.Type == html.ElementNode && child.Data == "br" {
			text.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(text.String()), " ")
}
