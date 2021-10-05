package gits

import (
	"bufio"
	"bytes"
	"github.com/gookit/color"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Status struct {
	workplace string
	branch    string
	commit    string
	upstream  string
	ahead     int
	behind    int
	unTracked int
	unmerged  int
	unStaged  gitArea
	staged    gitArea
	newPull   bool
	tag       string
}

func NewStatus(workplace string) *Status {
	return &Status{workplace: workplace}
}

func (status *Status) Parse(b []byte) error {
	var r = bytes.NewReader(b)

	var err error
	var s = bufio.NewScanner(r)

	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}

		status.parseLine(s.Text())
	}

	return err
}

func (status *Status) Fmt() string {
	var (
		modifiedGlyph  = "Δ"
		dirtyGlyph     = "✘"
		cleanGlyph     = "✔"
		unTrackedGlyph = "?"
		unmergedGlyph  = "‼"
		aheadArrow     = "↑"
		behindArrow    = "↓"
		newPullGlyph   = "🔥"
	)

	branchFmt := color.New(color.FgBlue)

	aheadFmt := color.New(color.OpStrikethrough, color.BgYellow, color.FgBlack)
	behindFmt := color.New(color.OpStrikethrough, color.BgRed, color.FgWhite)

	modifiedFmt := color.New(color.FgRed)
	dirtyFmt := color.New(color.FgRed)
	cleanFmt := color.New(color.FgGreen)

	unTrackedFmt := color.New(color.OpStrikethrough)
	unmergedFmt := color.New(color.FgCyan)

	newPullFmt := color.New(color.FgRed)

	var buf bytes.Buffer
	buf.WriteString(branchFmt.Sprint(status.branch))
	buf.WriteRune(' ')

	if status.tag != "" {
		buf.WriteString(color.New(color.FgYellow).Sprint(status.tag))
		buf.WriteRune(' ')
	}

	if status.ahead > 0 {
		buf.WriteString(aheadFmt.Sprint(" ", aheadArrow, status.ahead, " "))
	}
	if status.behind > 0 {
		buf.WriteString(behindFmt.Sprint(" ", behindArrow, status.behind, " "))
	}
	if status.ahead > 0 || status.behind > 0 {
		buf.WriteRune(' ')
	}
	if status.unTracked > 0 {
		buf.WriteString(unTrackedFmt.Sprint(unTrackedGlyph))
		buf.WriteRune(' ')
	}
	if status.hasUnmerged() {
		buf.WriteString(unmergedFmt.Sprint(unmergedGlyph))
		buf.WriteRune(' ')
	}
	if status.hasModified() {
		buf.WriteString(modifiedFmt.Sprint(modifiedGlyph))
		buf.WriteRune(' ')
	}
	if status.IsDirty() {
		buf.WriteString(dirtyFmt.Sprint(dirtyGlyph))
	} else {
		buf.WriteString(cleanFmt.Sprint(cleanGlyph))
	}
	if status.newPull {
		buf.WriteRune(' ')
		buf.WriteString(newPullFmt.Sprintf(newPullGlyph))
	}
	return buf.String()
}

func (status *Status) CanPull(force bool) bool {
	if status.behind <= 0 {
		return false
	}
	if force {
		return true
	}
	if status.ahead > 0 {
		return false
	}
	if status.IsDirty() {
		return false
	}
	if status.unTracked > 0 {
		return false
	}
	if status.hasModified() {
		return false
	}
	return true
}

func (status *Status) CanPush(force bool) bool {
	if status.ahead <= 0 {
		return false
	}
	if force {
		return true
	}
	if status.behind > 0 {
		return false
	}
	if status.IsDirty() {
		return false
	}
	return true
}

func (status *Status) IsDirty() bool {
	return status.staged.hasChanged()
}

func (status *Status) Branch() string {
	return status.branch
}

func (status *Status) SetNewPull() {
	status.newPull = true
}

func (status *Status) SetTag(tag string) {
	status.tag = tag
}

func (status *Status) hasUnmerged() bool {
	if status.unmerged > 0 {
		return true
	}
	gitDir, err := pathToGitDir(status.workplace)
	if err != nil {
		return false
	}

	_, err = ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD"))
	return err == nil
}

func (status *Status) hasModified() bool {
	return status.unStaged.hasChanged()
}

func (status *Status) parseLine(line string) {
	s := bufio.NewScanner(strings.NewReader(line))
	s.Split(bufio.ScanWords)

	for s.Scan() {
		switch s.Text() {
		case "#":
			_ = status.parseBranchInfo(s)
		case "1":
			_ = status.parseTrackedFile(s)
		case "2":
			_ = status.parseRenamedFile(s)
		case "u":
			status.unmerged++
		case "?":
			status.unTracked++
		}
	}
}

func (status *Status) parseBranchInfo(s *bufio.Scanner) (err error) {
	for s.Scan() {
		switch s.Text() {
		case "branch.oid":
			status.commit = consumeNext(s)
		case "branch.head":
			status.branch = consumeNext(s)
		case "branch.upstream":
			status.upstream = consumeNext(s)
		case "branch.ab":
			err = status.parseAheadBehind(s)
		}
	}
	return err
}

func (status *Status) parseAheadBehind(s *bufio.Scanner) error {
	for s.Scan() {
		i, err := strconv.Atoi(s.Text()[1:])
		if err != nil {
			return err
		}

		switch s.Text()[:1] {
		case "+":
			status.ahead = i
		case "-":
			status.behind = i
		}
	}
	return nil
}

func (status *Status) parseTrackedFile(s *bufio.Scanner) error {
	var index int
	for s.Scan() {
		switch index {
		case 0: // xy
			status.parseXY(s.Text())
		default:
			continue
		}
		index++
	}
	return nil
}

func (status *Status) parseXY(xy string) {
	switch xy[:1] {
	case "M":
		status.staged.modified++
	case "A":
		status.staged.added++
	case "D":
		status.staged.deleted++
	case "R":
		status.staged.renamed++
	case "C":
		status.staged.copied++
	}

	switch xy[1:] {
	case "M":
		status.unStaged.modified++
	case "A":
		status.unStaged.added++
	case "D":
		status.unStaged.deleted++
	case "R":
		status.unStaged.renamed++
	case "C":
		status.unStaged.copied++
	}
}

func (status *Status) parseRenamedFile(s *bufio.Scanner) error {
	return status.parseTrackedFile(s)
}

type gitArea struct {
	modified int
	added    int
	deleted  int
	renamed  int
	copied   int
}

func (a *gitArea) hasChanged() bool {
	var changed bool
	if a.added != 0 {
		changed = true
	}
	if a.deleted != 0 {
		changed = true
	}
	if a.modified != 0 {
		changed = true
	}
	if a.copied != 0 {
		changed = true
	}
	if a.renamed != 0 {
		changed = true
	}
	return changed
}

func consumeNext(s *bufio.Scanner) string {
	if s.Scan() {
		return s.Text()
	}
	return ""
}

func pathToGitDir(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Dir = cwd

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
