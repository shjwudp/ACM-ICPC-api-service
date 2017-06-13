package model

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// Problem part of xml
type Problem struct {
	ID               int    `xml:"id,attr"`
	NumberSolved     int    `xml:"numberSolved,attr"`
	BestSolutionTime int    `xml:"bestSolutionTime,attr"`
	LastSolutionTime int    `xml:"lastSolutionTime,attr"`
	Title            string `xml:"title,attr"`
	ShortName        string
}

// StandingsHeader part of xml
type StandingsHeader struct {
	CurrentTimeStamp  int64
	CurrentDate       string     `xml:"currentDate,attr"`
	ProblemCount      int        `xml:"problemCount,attr"`
	ProblemsAttempted int        `xml:"problemsAttempted,attr"`
	TotalAttempts     int        `xml:"totalAttempts,attr"`
	TotalSolved       int        `xml:"totalSolved,attr"`
	Problems          []*Problem `xml:"problem"`
}

// ProblemSummaryInfo part of xml
type ProblemSummaryInfo struct {
	Index         int  `xml:"index,attr"`
	Attempts      int  `xml:"attempts,attr"`
	IsPending     bool `xml:"isPending,attr"`
	IsSolved      bool `xml:"isSolved,attr"`
	Points        int  `xml:"points,attr"`
	SolutionTime  int  `xml:"solutionTime,attr"`
	IsFirstSolved bool
}

// TeamStanding part of xml
type TeamStanding struct {
	FirstSolved         int                   `xml:"firstSolved,attr"`
	Index               int                   `xml:"index,attr"`
	LastSolved          int                   `xml:"lastSolved,attr"`
	Points              int                   `xml:"points,attr"`
	ProblemsAttempted   int                   `xml:"problemsAttempted,attr"`
	Rank                int                   `xml:"rank,attr"`
	Solved              int                   `xml:"solved,attr"`
	TeamName            string                `xml:"teamName,attr"`
	TeamKey             string                `xml:"teamKey,attr"`
	TotalAttempts       int                   `xml:"totalAttempts,attr"`
	ProblemSummaryInfos []*ProblemSummaryInfo `xml:"problemSummaryInfo"`
	NickName            string
	School              string
	IsStar              bool
	IsGirl              bool
	Coach               string
	Player1             string
	Player2             string
	Player3             string
	SeatID              string
}

// ContestStanding part of xml
type ContestStanding struct {
	XMLName         xml.Name        `xml:"contestStandings"`
	StandingsHeader StandingsHeader `xml:"standingsHeader"`
	TeamStandings   []*TeamStanding `xml:"teamStanding"`
}

type problemArray []*Problem

func (a problemArray) Len() int           { return len(a) }
func (a problemArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a problemArray) Less(i, j int) bool { return a[i].ID < a[j].ID }

type problemSummaryInfoArray []*ProblemSummaryInfo

func (a problemSummaryInfoArray) Len() int           { return len(a) }
func (a problemSummaryInfoArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a problemSummaryInfoArray) Less(i, j int) bool { return a[i].Index < a[j].Index }

type teamStandingArray []*TeamStanding

func (a teamStandingArray) Len() int           { return len(a) }
func (a teamStandingArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a teamStandingArray) Less(i, j int) bool { return a[i].Rank < a[j].Rank }

// ParseResultXML parse(xmlpath) => go struct
func ParseResultXML(resultXMLPath string) (*ContestStanding, error) {
	xmlFile, err := os.Open(resultXMLPath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}
	cs := new(ContestStanding)
	err = xml.Unmarshal([]byte(data), &cs)
	if err != nil {
		return nil, err
	}
	time, err := time.Parse(
		"Mon Jan 2 15:04:05 MST 2006",
		cs.StandingsHeader.CurrentDate,
	)
	if err != nil {
		return nil, err
	}
	cs.StandingsHeader.CurrentTimeStamp = time.Unix()
	sort.Sort(problemArray(cs.StandingsHeader.Problems))
	sort.Sort(teamStandingArray(cs.TeamStandings))
	problemList := cs.StandingsHeader.Problems
	for _, p := range problemList {
		tmp := p.ID
		for tmp > 0 {
			p.ShortName += string('A' + tmp%26)
			tmp /= 26
		}
	}
	problems := cs.StandingsHeader.Problems
	for _, t := range cs.TeamStandings {
		sort.Sort(problemSummaryInfoArray(t.ProblemSummaryInfos))
		for _, p := range t.ProblemSummaryInfos {
			if p.IsSolved && p.SolutionTime == problems[p.Index-1].BestSolutionTime {
				p.IsFirstSolved = true
			}
		}
	}
	return cs, nil
}
