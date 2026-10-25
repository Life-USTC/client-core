package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"time"
)

const (
	YJS1ScheduleAppURL  = "https://yjs1.ustc.edu.cn/gsapp/sys/kbcxappustc/*default/index.do?THEME=blue&EMAP_LANG=zh&min=1#/xskbcx"
	YJS1ScheduleTermURL = "https://yjs1.ustc.edu.cn/gsapp/sys/kbcxappustc/modules/xskbcx/xnxqxxcx.do"
	YJS1ScheduleRowsURL = "https://yjs1.ustc.edu.cn/gsapp/sys/kbcxappustc/modules/xskbcx/xskbxxcx.do"

	YJS1ScoreAppURL     = "https://yjs1.ustc.edu.cn/gsapp/sys/wdcjapp/*default/index.do?THEME=blue&EMAP_LANG=zh&min=1#/wdcj"
	YJS1ScoreInfoURL    = "https://yjs1.ustc.edu.cn/gsapp/sys/wdcjapp/modules/wdcj/queryInfo.do"
	YJS1ScoreProfileURL = "https://yjs1.ustc.edu.cn/gsapp/sys/wdcjapp/modules/wdcj/hqxh.do"
	YJS1ScoreRowsURL    = "https://yjs1.ustc.edu.cn/gsapp/sys/wdcjapp/modules/wdcj/xscjcx.do"

	YJS1ExamAppURL  = "https://yjs1.ustc.edu.cn/gsapp/sys/kssbappustc/*default/index.do?THEME=blue&EMAP_LANG=zh&min=1#/xskssbcx"
	YJS1ExamTermURL = "https://yjs1.ustc.edu.cn/gsapp/sys/kssbappustc/modules/xskssbcx/hqxnxq.do"
	YJS1ExamRowsURL = "https://yjs1.ustc.edu.cn/gsapp/sys/kssbappustc/modules/xskssbcx/xscxkssbxx.do"

	YJS1HomeworkAppURL  = "https://yjs1.ustc.edu.cn/gsapp/sys/wdzyappustc/*default/index.do?THEME=blue&EMAP_LANG=zh&min=1#/wdzy"
	YJS1HomeworkTermURL = "https://yjs1.ustc.edu.cn/gsapp/sys/wdzyappustc/modules/wdzy/hqxnxq.do"
	YJS1HomeworkRowsURL = "https://yjs1.ustc.edu.cn/gsapp/sys/wdzyappustc/modules/wdzy/wdzycx.do"
)

// Each form retains the fields used by the corresponding existing YJS1 app.
// Clients encode these values and retain HTTP headers, authentication and sessions.
func YJS1ScheduleTermsForm() url.Values {
	return url.Values{"SFSY": {"1"}, "pageSize": {"100"}, "pageNumber": {"1"}}
}

func YJS1ScheduleRowsForm(semesterCode string) url.Values {
	return url.Values{
		"querySetting": {yjs1QuerySetting([]yjs1Query{
			{Name: "XNXQDM", Builder: "equal", LinkOpt: "AND", BuilderList: "cbl_String", Value: semesterCode},
		})},
		"pageSize": {"999"}, "pageNumber": {"1"},
	}
}

func YJS1ExamRowsForm(semesterCode string) url.Values {
	return url.Values{
		"querySetting": {yjs1QuerySetting([]yjs1Query{
			{Name: "XNXQDM", Caption: "学年学期", Builder: "m_value_equal", LinkOpt: "AND", Value: semesterCode},
			{Name: "SBZT", Caption: "申报状态", Builder: "equal", LinkOpt: "AND", Value: "1"},
		})},
		"pageSize": {"999"}, "pageNumber": {"1"},
	}
}

func YJS1ScoreRowsForm() url.Values {
	return url.Values{"pageSize": {"999"}, "pageNumber": {"1"}}
}

func YJS1HomeworkTermsForm() url.Values {
	return url.Values{"*order": {"-XNDM,-XQDM"}}
}

func YJS1HomeworkRowsForm(semesterCode, userID string) url.Values {
	return url.Values{
		"XH": {userID},
		"querySetting": {yjs1QuerySetting([]yjs1Query{
			{Name: "XNXQDM", Builder: "m_value_equal", Value: semesterCode},
		})},
		"pageSize": {"999"}, "pageNumber": {"1"},
	}
}

type yjs1Query struct {
	Name        string `json:"name"`
	Caption     string `json:"caption,omitempty"`
	Builder     string `json:"builder,omitempty"`
	LinkOpt     string `json:"linkOpt,omitempty"`
	BuilderList string `json:"builderList,omitempty"`
	Value       string `json:"value,omitempty"`
}

func yjs1QuerySetting(filters []yjs1Query) string {
	data, _ := json.Marshal(filters)
	return string(data)
}

// parseYJS1Rows retains the existing first-JSON-value decoder and first data set
// behavior. It does not choose a data-set name or merge multiple data sets.
func parseYJS1Rows[T any](data []byte) ([]T, error) {
	var response yjs1RowsResponse[T]
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&response); err != nil {
		return nil, err
	}
	if response.Code != "" && response.Code != "0" {
		return nil, fmt.Errorf("graduate yjs1 request failed: %s", response.Msg)
	}
	for _, set := range response.Datas {
		return set.Rows, nil
	}
	return nil, nil
}

// ParseYJS1Terms preserves source order and the current-term marker. Selecting a
// requested or default semester remains the client's policy. Missing/null rows
// remain nil so the CLI retains its existing term-cache refetch behavior.
func ParseYJS1Terms(data []byte) ([]YJS1Semester, error) {
	rows, err := parseYJS1Rows[yjs1Term](data)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	return yjs1SemestersFromTerms(rows), nil
}

func ParseYJS1Curriculum(data []byte, semesterID int) ([]YJS1CurriculumItem, error) {
	rows, err := parseYJS1Rows[yjs1ScheduleRow](data)
	if err != nil {
		return nil, err
	}
	return yjs1CurriculumItemsFromRows(YJS1Semester{ID: semesterID}, rows), nil
}

// Score info and profile parsing remain separate so the client can preserve its
// existing request sequence and stop before subsequent requests on an error.
func ParseYJS1ScoreInfo(data []byte) ([]YJS1ScoreInfo, error) {
	return parseYJS1Rows[YJS1ScoreInfo](data)
}

func ParseYJS1ScoreProfile(data []byte) (YJS1ScoreProfile, error) {
	var profile YJS1ScoreProfile
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&profile); err != nil {
		return YJS1ScoreProfile{}, err
	}
	if profile.Code != "" && profile.Code != "0" {
		return YJS1ScoreProfile{}, fmt.Errorf("graduate score profile request failed: %s", profile.Msg)
	}
	return profile, nil
}

func ParseYJS1Scores(data []byte, info []YJS1ScoreInfo, profile YJS1ScoreProfile) (YJS1ScoreReport, error) {
	rows, err := parseYJS1Rows[yjs1ScoreRow](data)
	if err != nil {
		return YJS1ScoreReport{}, err
	}
	report := YJS1ScoreReport{Items: yjs1ScoreItemsFromRows(rows), Summary: buildYJS1ScoreSummary(firstYJS1ScoreInfoRow(info), profile)}
	sortYJS1Scores(report.Items)
	return report, nil
}

func ParseYJS1Exams(data []byte) ([]YJS1ExamItem, error) {
	rows, err := parseYJS1Rows[yjs1ExamRow](data)
	if err != nil {
		return nil, err
	}
	return yjs1ExamItemsFromRows(rows), nil
}

// ParseYJS1Homework requires the caller's current time and nonnil local timezone;
// it does not choose a timezone or read a clock. Status precedence is unchanged.
func ParseYJS1Homework(data []byte, now time.Time, location *time.Location) ([]YJS1HomeworkItem, error) {
	rows, err := parseYJS1Rows[yjs1HomeworkRow](data)
	if err != nil {
		return nil, err
	}
	return yjs1HomeworkItemsFromRows(rows, now, location), nil
}

var yjs1UserIDPattern = regexp.MustCompile(`"USERID":"([^"]+)"`)

func ParseYJS1UserID(data []byte, appPageURL string) (string, error) {
	match := yjs1UserIDPattern.FindSubmatch(data)
	if len(match) < 2 {
		return "", fmt.Errorf("could not resolve graduate student ID from %s", appPageURL)
	}
	return string(match[1]), nil
}
