package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	ScoreSemestersURL = "https://jw.ustc.edu.cn/for-std/grade/sheet/getSemesters"
	ScoreTypesURL     = "https://jw.ustc.edu.cn/for-std/grade/sheet/getGradeSheetTypes"
	scoreListURL      = "https://jw.ustc.edu.cn/for-std/grade/sheet/getGradeList"
)

// NumericScoreListURL retains the CLI's explicit semester selection.
func NumericScoreListURL(trainTypeID int, semesterIDs []int) string {
	ids := make([]string, len(semesterIDs))
	for i, id := range semesterIDs {
		ids[i] = strconv.Itoa(id)
	}
	query := url.Values{}
	query.Set("trainTypeId", strconv.Itoa(trainTypeID))
	query.Set("semesterIds", strings.Join(ids, ","))
	return scoreListURL + "?" + query.Encode()
}

// ScalarScoreListURL retains the iOS request for all score semesters.
func ScalarScoreListURL(trainTypeID int) string {
	return fmt.Sprintf("%s?trainTypeId=%d&semesterIds", scoreListURL, trainTypeID)
}

type ScoreOverview struct {
	PassedCredits    *float64 `json:"passedCredits"`
	GPA              *float64 `json:"gpa"`
	WeightedScore    *float64 `json:"weightedScore"`
	NotPassedCredits *float64 `json:"notPassedCredits"`
	ArithmeticScore  *float64 `json:"arithmeticScore"`
}

type ScoreRank struct {
	StartSemesterName string   `json:"startSemesterName"`
	EndSemesterName   string   `json:"endSemesterName"`
	GPA               *float64 `json:"gpa"`
	Grade             string   `json:"grade"`
	MajorName         *string  `json:"majorName"`
	MajorRank         *int     `json:"majorRank"`
	MajorStdCount     *int     `json:"majorStdCount"`
}

// Score preserves the scalar semester identifier and an absent grade point.
type Score struct {
	SemesterID   string   `json:"semesterId"`
	SemesterName string   `json:"semesterName"`
	CourseName   string   `json:"courseName"`
	LessonCode   string   `json:"lessonCode"`
	CourseCode   string   `json:"courseCode"`
	Credits      float64  `json:"credits"`
	GradePoint   *float64 `json:"gradePoint"`
	Score        string   `json:"score"`
	GradeText    string   `json:"gradeText"`
}

type ScoreReport struct {
	Overview ScoreOverview   `json:"overview"`
	Rank     *ScoreRank      `json:"rank"`
	Items    []Score         `json:"items"`
	Summary  json.RawMessage `json:"summary,omitempty"`
}

type numericScoreEntry struct {
	SemesterID    int     `json:"semesterId"`
	SemesterAssoc int     `json:"semesterAssoc"`
	SemesterName  string  `json:"semesterName"`
	SemesterCh    string  `json:"semesterCh"`
	SemesterEn    string  `json:"semesterEn"`
	LessonCode    string  `json:"lessonCode"`
	CourseCode    string  `json:"courseCode"`
	CourseName    string  `json:"courseName"`
	CourseNameCh  string  `json:"courseNameCh"`
	CourseNameEn  string  `json:"courseNameEn"`
	Credits       float64 `json:"credits"`
	GradePoint    float64 `json:"gp"`
	Score         string  `json:"score"`
	ScoreCh       string  `json:"scoreCh"`
	ScoreEn       string  `json:"scoreEn"`
	ScoreText     string  `json:"scoreText"`
	Grade         string  `json:"grade"`
	GradeText     string  `json:"gradeDetail"`
	Course        struct {
		Code    string  `json:"code"`
		NameZh  string  `json:"nameZh"`
		NameEn  string  `json:"nameEn"`
		Credits float64 `json:"credits"`
	} `json:"course"`
	Semester struct {
		ID         int    `json:"id"`
		Code       string `json:"code"`
		SemesterCn string `json:"semesterCn"`
		SemesterEn string `json:"semesterEn"`
		StartDate  string `json:"startDate"`
		EndDate    string `json:"endDate"`
		IsLast     bool   `json:"isLast"`
	} `json:"semester"`
}

// ParseNumericScores preserves the CLI's strict numeric response schema,
// field precedence, summary and semester/name ordering. Missing grade points
// remain zero in this profile, as in the existing CLI output.
func ParseNumericScores(data []byte, semesterNames map[int]string) (ScoreReport, error) {
	var payload struct {
		Overview     ScoreOverview `json:"overview"`
		StdGradeRank *ScoreRank    `json:"stdGradeRank"`
		Semesters    []struct {
			ID     int                 `json:"id"`
			Scores []numericScoreEntry `json:"scores"`
		} `json:"semesters"`
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&payload); err != nil {
		return ScoreReport{}, err
	}
	report := ScoreReport{Overview: payload.Overview, Rank: payload.StdGradeRank, Items: []Score{}}
	for _, semester := range payload.Semesters {
		for _, item := range semester.Scores {
			semesterID := item.SemesterID
			if semesterID == 0 {
				semesterID = item.SemesterAssoc
			}
			if semesterID == 0 {
				semesterID = semester.ID
			}
			credits := item.Credits
			if credits == 0 {
				credits = item.Course.Credits
			}
			gradePoint := item.GradePoint
			report.Items = append(report.Items, Score{
				SemesterID:   strconv.Itoa(semesterID),
				SemesterName: firstScoreText(item.SemesterName, item.SemesterCh, semesterNames[semester.ID], item.Semester.SemesterCn, item.Semester.SemesterEn),
				CourseName:   firstScoreText(item.CourseName, item.CourseNameCh, item.Course.NameZh, item.Course.NameEn, item.CourseNameEn),
				LessonCode:   item.LessonCode, CourseCode: firstScoreText(item.CourseCode, item.Course.Code),
				Credits: credits, GradePoint: &gradePoint,
				Score:     firstScoreText(item.ScoreCh, item.Score, item.ScoreText, item.Grade, item.ScoreEn),
				GradeText: firstScoreText(item.GradeText, item.ScoreText),
			})
		}
	}
	sort.Slice(report.Items, func(i, j int) bool {
		a, b := report.Items[i], report.Items[j]
		aID, _ := strconv.Atoi(a.SemesterID)
		bID, _ := strconv.Atoi(b.SemesterID)
		if aID != bID {
			return aID > bID
		}
		if a.CourseName != b.CourseName {
			return a.CourseName < b.CourseName
		}
		return a.LessonCode < b.LessonCode
	})
	var err error
	report.Summary, err = numericScoreSummary(report.Overview, report.Rank)
	return report, err
}

// ParseScalarScores preserves the iOS helper's SwiftyJSON scalar conversions
// and response order. Metadata stays optional so the native adapter retains
// its existing display defaults, including the missing-major-name label.
func ParseScalarScores(data []byte) (ScoreReport, error) {
	value, err := decodeScalarJSON(data)
	if err != nil {
		return ScoreReport{}, err
	}
	root := scalarJSONObject(value)
	overview := scalarJSONObject(root["overview"])
	report := ScoreReport{Overview: ScoreOverview{GPA: scalarJSONNumber(overview["gpa"])}, Items: []Score{}}
	if rank, ok := root["stdGradeRank"].(map[string]any); ok {
		report.Rank = &ScoreRank{
			MajorName: scalarJSONString(rank["majorName"]),
			MajorRank: scalarJSONInt(rank["majorRank"]), MajorStdCount: scalarJSONInt(rank["majorStdCount"]),
		}
	}
	for _, semester := range scalarJSONChildren(root["semesters"]) {
		for _, entry := range scalarJSONChildren(scalarJSONObject(semester)["scores"]) {
			item := scalarJSONObject(entry)
			credits, _ := strconv.ParseFloat(scalarJSONText(item["credits"]), 64)
			var gradePoint *float64
			if number, err := strconv.ParseFloat(scalarJSONText(item["gp"]), 64); err == nil {
				gradePoint = &number
			}
			report.Items = append(report.Items, Score{
				SemesterID: scalarJSONText(item["semesterAssoc"]), SemesterName: scalarJSONText(item["semesterCh"]),
				CourseName: scalarJSONText(item["courseNameCh"]), CourseCode: scalarJSONText(item["courseCode"]),
				LessonCode: scalarJSONText(item["lessonCode"]), Credits: credits, GradePoint: gradePoint,
				Score: scalarJSONText(item["scoreCh"]),
			})
		}
	}
	return report, nil
}

func numericScoreSummary(overview ScoreOverview, rank *ScoreRank) (json.RawMessage, error) {
	summary := map[string]any{}
	passed, failed, gpa := scoreValue(overview.PassedCredits), scoreValue(overview.NotPassedCredits), scoreValue(overview.GPA)
	weighted, arithmetic := scoreValue(overview.WeightedScore), scoreValue(overview.ArithmeticScore)
	if passed != 0 || failed != 0 || gpa != 0 || weighted != 0 || arithmetic != 0 {
		summary["totalCredits"], summary["earnedCredits"], summary["failedCredits"] = passed+failed, passed, failed
		summary["gpa"], summary["weightedAverage"], summary["averageScore"] = gpa, weighted, arithmetic
	}
	if rank != nil {
		summary["fromSemester"], summary["toSemester"], summary["periodGPA"] = rank.StartSemesterName, rank.EndSemesterName, scoreValue(rank.GPA)
		if rank.MajorName != nil && *rank.MajorName != "" && rank.MajorRank != nil && *rank.MajorRank != 0 && rank.MajorStdCount != nil && *rank.MajorStdCount != 0 {
			summary["ranking"] = fmt.Sprintf("全校%s年级%s专业 GPA 排名 %d/%d", rank.Grade, *rank.MajorName, *rank.MajorRank, *rank.MajorStdCount)
		}
	}
	if len(summary) == 0 {
		return nil, nil
	}
	return json.Marshal(summary)
}

func scoreValue(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func firstScoreText(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
