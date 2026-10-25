package ustc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type yjs1RowsResponse[T any] struct {
	Code  string                            `json:"code"`
	Msg   string                            `json:"msg"`
	Datas map[string]yjs1RowsResponseSet[T] `json:"datas"`
}

type yjs1RowsResponseSet[T any] struct {
	Rows []T `json:"rows"`
}

type yjs1Term struct {
	DM     string `json:"DM"`
	MC     string `json:"MC"`
	SFDQXQ string `json:"SFDQXQ"`
	QSSJ   string `json:"QSSJ"`
	JZSJ   string `json:"JZSJ"`
}

type yjs1ScheduleRow struct {
	WID    string `json:"WID"`
	BJMC   string `json:"BJMC"`
	KCDM   string `json:"KCDM"`
	KCMC   string `json:"KCMC"`
	PKSJDD string `json:"PKSJDD"`
	RKJS   string `json:"RKJS"`
	ZCMC   string `json:"ZCMC"`
}

// YJS1ScoreInfo is the profile fallback returned before the detailed score profile.
type YJS1ScoreInfo struct {
	XH            string `json:"XH"`
	PYCCDMDisplay string `json:"PYCCDM_DISPLAY"`
	ZYDMDisplay   string `json:"ZYDM_DISPLAY"`
	YXDMDisplay   string `json:"YXDM_DISPLAY"`
	DSXM          string `json:"DSXM"`
	NJDMDisplay   string `json:"NJDM_DISPLAY"`
	ZCZTDisplay   string `json:"ZCZT_DISPLAY"`
	YJBYSJ        string `json:"YJBYSJ"`
}

// YJS1ScoreProfile supplies preferred summary fields from the hqxh endpoint.
type YJS1ScoreProfile struct {
	Code   string                  `json:"code"`
	Msg    string                  `json:"msg"`
	XH     string                  `json:"XH"`
	XSCJXX YJS1ScoreProfileDetails `json:"xscjxx"`
}

type YJS1ScoreProfileDetails struct {
	PYCCDMDisplay string `json:"PYCCDM_DISPLAY"`
	ZYDMDisplay   string `json:"ZYDM_DISPLAY"`
	YXDMDisplay   string `json:"YXDM_DISPLAY"`
	DSXM          string `json:"DSXM"`
	NJDMDisplay   string `json:"NJDM_DISPLAY"`
	ZCZTDisplay   string `json:"ZCZT_DISPLAY"`
	YJBYSJ        string `json:"YJBYSJ"`
	PJJDZ         string `json:"PJJDZ"`
}

type yjs1ScoreRow struct {
	XNXQDM        string `json:"XNXQDM"`
	XNXQDMDisplay string `json:"XNXQDM_DISPLAY"`
	BJMC          string `json:"BJMC"`
	KCDM          string `json:"KCDM"`
	KCMC          string `json:"KCMC"`
	CJFZDMDisplay string `json:"CJFZDM_DISPLAY"`
	CJJLDisplay   string `json:"CJJL_DISPLAY"`
	XF            any    `json:"XF"`
	CJ            any    `json:"CJ"`
	JDZ           any    `json:"JDZ"`
}

type yjs1ExamRow struct {
	BJMC        string `json:"BJMC"`
	KCDM        string `json:"KCDM"`
	KCMC        string `json:"KCMC"`
	KSDD        string `json:"KSDD"`
	KSKSSJ      string `json:"KSKSSJ"`
	KSJSSJ      string `json:"KSJSSJ"`
	KSXS        string `json:"KSXS"`
	KSXSDisplay string `json:"KSXS_DISPLAY"`
	SBZT        string `json:"SBZT"`
	SBZTDisplay string `json:"SBZT_DISPLAY"`
	ZWH         string `json:"ZWH"`
}

type yjs1HomeworkRow struct {
	WID     string `json:"WID"`
	XSZYWID string `json:"XSZYWID"`
	KCMC    string `json:"KCMC"`
	ZYMC    string `json:"ZYMC"`
	ZYKSSJ  string `json:"ZYKSSJ"`
	ZYJZSJ  string `json:"ZYJZSJ"`
	ZYSCSJ  string `json:"ZYSCSJ"`
	CJ      any    `json:"CJ"`
}

// YJS1Semester retains the graduate term code, display names and current flag.
type YJS1Semester struct {
	ID         int    `json:"id"`
	Code       string `json:"code,omitempty"`
	SemesterCn string `json:"semesterCn,omitempty"`
	SemesterEn string `json:"semesterEn,omitempty"`
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
	IsLast     bool   `json:"isLast,omitempty"`
}

type YJS1CurriculumItem struct {
	SemesterID int      `json:"semesterId"`
	LessonID   int      `json:"lessonId,omitempty"`
	LessonCode string   `json:"lessonCode,omitempty"`
	CourseCode string   `json:"courseCode,omitempty"`
	CourseName string   `json:"courseName,omitempty"`
	Credits    float64  `json:"credits,omitempty"`
	Teachers   []string `json:"teachers,omitempty"`
	Schedule   string   `json:"schedule,omitempty"`
}

type YJS1ExamItem struct {
	CourseName string `json:"courseName,omitempty"`
	LessonCode string `json:"lessonCode,omitempty"`
	ExamType   string `json:"examType,omitempty"`
	DateTime   string `json:"dateTime,omitempty"`
	Location   string `json:"location,omitempty"`
	Seat       string `json:"seat,omitempty"`
	Status     string `json:"status,omitempty"`
}

type YJS1ScoreReport struct {
	Summary json.RawMessage `json:"summary,omitempty"`
	Items   []YJS1ScoreItem `json:"items"`
}

type YJS1ScoreItem struct {
	SemesterID   int     `json:"semesterId"`
	SemesterName string  `json:"semesterName,omitempty"`
	CourseName   string  `json:"courseName,omitempty"`
	LessonCode   string  `json:"lessonCode,omitempty"`
	CourseCode   string  `json:"courseCode,omitempty"`
	Credits      float64 `json:"credits,omitempty"`
	GradePoint   float64 `json:"gradePoint,omitempty"`
	Score        string  `json:"score,omitempty"`
	GradeText    string  `json:"gradeText,omitempty"`
}

type YJS1HomeworkItem struct {
	ID         string `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	CourseName string `json:"courseName,omitempty"`
	StartAt    string `json:"startAt,omitempty"`
	EndAt      string `json:"endAt,omitempty"`
	Status     string `json:"status,omitempty"`
}

func yjs1SemestersFromTerms(rows []yjs1Term) []YJS1Semester {
	semesters := make([]YJS1Semester, 0, len(rows))
	for _, row := range rows {
		id, _ := strconv.Atoi(row.DM)
		name := firstScoreText(row.MC, yjs1SemesterNameFromCode(row.DM))
		semesters = append(semesters, YJS1Semester{
			ID:         id,
			Code:       row.DM,
			SemesterCn: name,
			SemesterEn: name,
			StartDate:  row.QSSJ,
			EndDate:    row.JZSJ,
			IsLast:     row.SFDQXQ == "1",
		})
	}
	return semesters
}

func yjs1SemesterNameFromCode(code string) string {
	if len(code) < 5 {
		return code
	}
	year, err := strconv.Atoi(code[:4])
	if err != nil {
		return code
	}
	switch code[4:] {
	case "1":
		return fmt.Sprintf("%d年秋季学期", year)
	case "2":
		return fmt.Sprintf("%d年春季学期", year+1)
	case "3":
		return fmt.Sprintf("%d年夏季学期", year+1)
	default:
		return code
	}
}

func yjs1CurriculumItemsFromRows(semester YJS1Semester, rows []yjs1ScheduleRow) []YJS1CurriculumItem {
	items := make([]YJS1CurriculumItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, YJS1CurriculumItem{
			LessonID:   yjs1ParseLeadingInt(row.WID),
			SemesterID: semester.ID,
			LessonCode: yjs1CompactWhitespace(firstScoreText(row.BJMC, row.WID)),
			CourseCode: yjs1CompactWhitespace(row.KCDM),
			CourseName: yjs1CompactWhitespace(row.KCMC),
			Teachers:   yjs1SplitTeachers(row.RKJS),
			Schedule:   yjs1ScheduleString(row.PKSJDD, row.ZCMC),
		})
	}
	sortYJS1Curriculum(items)
	return items
}

func yjs1ScheduleString(lessonSchedule, weekRanges string) string {
	lessonSchedule = yjs1CompactWhitespace(lessonSchedule)
	weekRanges = yjs1CompactWhitespace(strings.ReplaceAll(weekRanges, "~", "-"))
	if lessonSchedule == "" {
		return weekRanges
	}
	if weekRanges == "" {
		return lessonSchedule
	}

	scheduleParts := splitSemicolonParts(lessonSchedule)
	weekParts := splitSemicolonParts(weekRanges)
	if len(scheduleParts) == len(weekParts) {
		parts := make([]string, 0, len(scheduleParts))
		for i, schedule := range scheduleParts {
			if week := weekParts[i]; week != "" {
				parts = append(parts, schedule+" ["+week+"周]")
				continue
			}
			parts = append(parts, schedule)
		}
		return strings.Join(parts, "; ")
	}

	return lessonSchedule + " [" + weekRanges + "周]"
}

func yjs1ScoreItemsFromRows(rows []yjs1ScoreRow) []YJS1ScoreItem {
	items := make([]YJS1ScoreItem, 0, len(rows))
	for _, row := range rows {
		semesterID, _ := strconv.Atoi(row.XNXQDM)
		items = append(items, YJS1ScoreItem{
			SemesterID:   semesterID,
			SemesterName: yjs1CompactWhitespace(firstScoreText(row.XNXQDMDisplay, yjs1SemesterNameFromCode(row.XNXQDM))),
			CourseName:   yjs1CompactWhitespace(row.KCMC),
			LessonCode:   yjs1CompactWhitespace(firstScoreText(row.BJMC, row.KCDM)),
			CourseCode:   yjs1CompactWhitespace(row.KCDM),
			Credits:      yjs1Float(row.XF),
			GradePoint:   yjs1Float(row.JDZ),
			Score:        yjs1String(row.CJ),
			GradeText:    yjs1Join(" / ", row.CJFZDMDisplay, row.CJJLDisplay),
		})
	}
	return items
}

func firstYJS1ScoreInfoRow(rows []YJS1ScoreInfo) YJS1ScoreInfo {
	if len(rows) == 0 {
		return YJS1ScoreInfo{}
	}
	return rows[0]
}

func buildYJS1ScoreSummary(info YJS1ScoreInfo, profile YJS1ScoreProfile) json.RawMessage {
	summary := map[string]any{}
	if studentID := yjs1CompactWhitespace(firstScoreText(profile.XH, info.XH)); studentID != "" {
		summary["studentId"] = studentID
	}
	if program := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.PYCCDMDisplay, info.PYCCDMDisplay)); program != "" {
		summary["program"] = program
	}
	if major := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.ZYDMDisplay, info.ZYDMDisplay)); major != "" {
		summary["major"] = major
	}
	if department := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.YXDMDisplay, info.YXDMDisplay)); department != "" {
		summary["department"] = department
	}
	if advisor := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.DSXM, info.DSXM)); advisor != "" {
		summary["advisor"] = advisor
	}
	if grade := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.NJDMDisplay, info.NJDMDisplay)); grade != "" {
		summary["grade"] = grade
	}
	if status := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.ZCZTDisplay, info.ZCZTDisplay)); status != "" {
		summary["status"] = status
	}
	if graduation := yjs1CompactWhitespace(firstScoreText(profile.XSCJXX.YJBYSJ, info.YJBYSJ)); graduation != "" {
		summary["expectedGraduation"] = graduation
	}
	if gpa := strings.TrimSpace(profile.XSCJXX.PJJDZ); gpa != "" {
		if value, err := strconv.ParseFloat(gpa, 64); err == nil {
			summary["averageGradePoint"] = value
		} else {
			summary["averageGradePoint"] = gpa
		}
	}
	if len(summary) == 0 {
		return nil
	}
	data, err := json.Marshal(summary)
	if err != nil {
		return nil
	}
	return data
}

func yjs1ExamItemsFromRows(rows []yjs1ExamRow) []YJS1ExamItem {
	items := make([]YJS1ExamItem, 0, len(rows))
	for _, row := range rows {
		dateTime := yjs1CompactWhitespace(strings.Trim(strings.Join([]string{
			strings.TrimSpace(row.KSKSSJ),
			strings.TrimSpace(row.KSJSSJ),
		}, " - "), " -"))
		items = append(items, YJS1ExamItem{
			CourseName: yjs1CompactWhitespace(row.KCMC),
			LessonCode: yjs1CompactWhitespace(firstScoreText(row.BJMC, row.KCDM)),
			DateTime:   dateTime,
			Location:   yjs1CompactWhitespace(row.KSDD),
			ExamType:   yjs1CompactWhitespace(firstScoreText(row.KSXSDisplay, row.KSXS)),
			Status:     yjs1CompactWhitespace(firstScoreText(row.SBZTDisplay, row.SBZT)),
			Seat:       yjs1CompactWhitespace(row.ZWH),
		})
	}
	sortYJS1Exams(items)
	return items
}

func yjs1HomeworkItemsFromRows(rows []yjs1HomeworkRow, now time.Time, location *time.Location) []YJS1HomeworkItem {
	items := make([]YJS1HomeworkItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, YJS1HomeworkItem{
			ID:         firstScoreText(row.XSZYWID, row.WID),
			Title:      yjs1CompactWhitespace(row.ZYMC),
			CourseName: yjs1CompactWhitespace(row.KCMC),
			StartAt:    yjs1CompactWhitespace(row.ZYKSSJ),
			EndAt:      yjs1CompactWhitespace(row.ZYJZSJ),
			Status:     yjs1HomeworkStatus(row, now, location),
		})
	}
	sortYJS1Homework(items)
	return items
}

func yjs1HomeworkStatus(row yjs1HomeworkRow, now time.Time, location *time.Location) string {
	if yjs1String(row.CJ) != "" {
		return "graded"
	}
	if strings.TrimSpace(row.XSZYWID) != "" || strings.TrimSpace(row.ZYSCSJ) != "" {
		return "submitted"
	}
	if deadline, ok := yjs1ParseTime(row.ZYJZSJ, location); ok && now.After(deadline) {
		return "overdue"
	}
	if start, ok := yjs1ParseTime(row.ZYKSSJ, location); ok && now.Before(start) {
		return "not started"
	}
	return "pending"
}

func yjs1ParseTime(value string, location *time.Location) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02"} {
		if parsed, err := time.ParseInLocation(layout, value, location); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func yjs1String(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return yjs1CompactWhitespace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case json.Number:
		return v.String()
	default:
		return yjs1CompactWhitespace(fmt.Sprint(v))
	}
}

func yjs1Float(value any) float64 {
	switch v := value.(type) {
	case nil:
		return 0
	case float64:
		return v
	case json.Number:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return f
	default:
		f, _ := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(v)), 64)
		return f
	}
}

func yjs1Join(sep string, values ...string) string {
	parts := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = yjs1CompactWhitespace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		parts = append(parts, value)
	}
	return strings.Join(parts, sep)
}

func splitSemicolonParts(value string) []string {
	rawParts := strings.Split(value, ";")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		part = yjs1CompactWhitespace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func sortYJS1Curriculum(items []YJS1CurriculumItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].LessonCode != items[j].LessonCode {
			return items[i].LessonCode < items[j].LessonCode
		}
		if items[i].CourseName != items[j].CourseName {
			return items[i].CourseName < items[j].CourseName
		}
		return items[i].LessonID < items[j].LessonID
	})
}

func sortYJS1Exams(items []YJS1ExamItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].DateTime != items[j].DateTime {
			return items[i].DateTime < items[j].DateTime
		}
		if items[i].CourseName != items[j].CourseName {
			return items[i].CourseName < items[j].CourseName
		}
		return items[i].LessonCode < items[j].LessonCode
	})
}

func sortYJS1Scores(items []YJS1ScoreItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].SemesterID != items[j].SemesterID {
			return items[i].SemesterID > items[j].SemesterID
		}
		if items[i].CourseName != items[j].CourseName {
			return items[i].CourseName < items[j].CourseName
		}
		return items[i].LessonCode < items[j].LessonCode
	})
}

func sortYJS1Homework(items []YJS1HomeworkItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].EndAt != items[j].EndAt {
			return items[i].EndAt < items[j].EndAt
		}
		if items[i].CourseName != items[j].CourseName {
			return items[i].CourseName < items[j].CourseName
		}
		return items[i].Title < items[j].Title
	})
}

func yjs1SplitTeachers(value string) []string {
	parts := regexp.MustCompile(`[、,，;/；]+`).Split(value, -1)
	teachers := make([]string, 0, len(parts))
	for _, part := range parts {
		name := yjs1CompactWhitespace(part)
		if name != "" {
			teachers = append(teachers, name)
		}
	}
	slices.Sort(teachers)
	return slices.Compact(teachers)
}

func yjs1ParseLeadingInt(value string) int {
	match := regexp.MustCompile(`^\d+`).FindString(strings.TrimSpace(value))
	if match == "" {
		return 0
	}
	parsed, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return parsed
}

func yjs1CompactWhitespace(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}
