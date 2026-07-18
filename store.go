package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type MemoCommand interface {
	AddMemo(stdout io.Writer, stderr io.Writer, args []string) int
	ListMemos(stdout io.Writer, stderr io.Writer, args []string) int
	ShowMemo(stdout io.Writer, stderr io.Writer, args []string) int
	SearchMemos(stdout io.Writer, stderr io.Writer, args []string) int
	DeleteMemo(stdout io.Writer, stderr io.Writer, args []string) int
}

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct {
}

func (rtp *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type MemoCommandImpl struct {
	Command      MemoCommand
	MemoPath     string
	TimeProvider TimeProvider
	MemoOperator MemoOperator
}

func (cmd *MemoCommandImpl) AddMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  memo add <title>")
		return 1
	}

	if isEmpty(args[0]) {
		printEmptyError(stderr, "タイトル")
		return 1
	}

	memos, err := cmd.MemoOperator.LoadMemos(cmd.MemoPath)
	if err != nil {
		printOpenFileError(stderr, err)
		return 1
	}

	now := cmd.TimeProvider.Now()

	addedMemo := createMemo(memos, args[0], now)

	memos = append(memos, addedMemo)

	if err := cmd.MemoOperator.SaveMemos(cmd.MemoPath, memos); err != nil {
		fmt.Fprintln(stderr, "メモの保存に失敗しました", err)
		return 1
	}

	return 0
}

func (cmd *MemoCommandImpl) ListMemos(stdout io.Writer, stderr io.Writer, args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  memo list")
		return 1
	}

	memos, err := cmd.MemoOperator.LoadMemos(cmd.MemoPath)
	if err != nil {
		printOpenFileError(stderr, err)
		return 1
	}

	if len(memos) == 0 {
		fmt.Fprintln(stdout, "No memos found.")
	} else {
		for _, memo := range memos {
			printMemoForList(stdout, memo)
		}
	}
	return 0
}

const DATE_TIME_FORMAT = "2006-01-02 15:04"

func (cmd *MemoCommandImpl) ShowMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  memo show <id>")
		return 1
	}

	id, err := resolveId(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	memos, err := cmd.MemoOperator.LoadMemos(cmd.MemoPath)
	if err != nil {
		printOpenFileError(stderr, err)
		return 1
	}

	memo, err := findById(memos, id)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	fmt.Fprintln(stdout, "# "+memo.Title)
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "ID: "+strconv.Itoa(memo.ID))
	fmt.Fprintln(stdout, "Created: "+memo.CreatedAt.Format(DATE_TIME_FORMAT))
	fmt.Fprintln(stdout, "Updated: "+memo.UpdatedAt.Format(DATE_TIME_FORMAT))
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, memo.Body)
	return 0
}

func (cmd *MemoCommandImpl) SearchMemos(stdout io.Writer, stderr io.Writer, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  memo search <keyword>")
		return 1
	}
	keyword := args[0]
	if isEmpty(keyword) {
		printEmptyError(stderr, "キーワード")
		return 1
	}

	memos, err := cmd.MemoOperator.LoadMemos(cmd.MemoPath)
	if err != nil {
		printOpenFileError(stderr, err)
		return 1
	}

	filteredMemos, err := filterMemos(memos, keyword)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	for _, memo := range filteredMemos {
		printMemoForList(stdout, memo)
	}

	return 0
}

func (cmd *MemoCommandImpl) DeleteMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  memo delete <id>")
		return 1
	}
	var id int
	id, err := resolveId(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	memos, err := cmd.MemoOperator.LoadMemos(cmd.MemoPath)
	if err != nil {
		printOpenFileError(stderr, err)
		return 1
	}

	newMemos, err := buildDeletedMemos(memos, id)

	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if err := cmd.MemoOperator.SaveMemos(cmd.MemoPath, newMemos); err != nil {
		fmt.Fprintln(stderr, "メモを削除できませんでした", err)
		return 1
	}
	return 0
}

func printOpenFileError(stderr io.Writer, err error) {
	fmt.Fprintln(stderr, "メモを開くことができませんでした", err)
}

func printMemoForList(stdout io.Writer, memo Memo) {
	fmt.Fprintf(stdout, "%d %s\t%s\n", memo.ID, memo.Title, memo.CreatedAt.Format("2006-01-02"))
}

func newNotFoundError(id int) error {
	return errors.New("memo not found: " + strconv.Itoa(id))
}

func createMemo(memos []Memo, title string, now time.Time) Memo {
	var id int
	if len(memos) > 0 {
		maxId := 0
		// JSONは順序を保証しないので全体を走査して最大IDを取得する
		for _, memo := range memos {
			if maxId < memo.ID {
				maxId = memo.ID
			}
		}
		id = maxId + 1
	} else {
		id = 1
	}

	return Memo{
		ID:        id,
		Title:     title,
		Body:      "",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func findById(memos []Memo, id int) (Memo, error) {

	for _, memo := range memos {
		if memo.ID == id {
			return memo, nil
		}
	}

	return Memo{}, newNotFoundError(id)
}

func filterMemos(memos []Memo, keyword string) ([]Memo, error) {
	filteredMemos := []Memo{}

	for _, memo := range memos {
		if strings.Contains(memo.Title, keyword) || strings.Contains(memo.Body, keyword) {
			filteredMemos = append(filteredMemos, memo)
		}
	}

	if len(filteredMemos) == 0 {
		return nil, errors.New("No matching memos found.")
	}

	return filteredMemos, nil
}

func buildDeletedMemos(memos []Memo, id int) ([]Memo, error) {

	deletedMemos := []Memo{}

	for _, memo := range memos {
		if memo.ID != id {
			deletedMemos = append(deletedMemos, memo)
		}
	}

	if len(deletedMemos) == len(memos) {
		return nil, newNotFoundError(id)
	}

	return deletedMemos, nil
}

func printEmptyError(stderr io.Writer, name string) {
	fmt.Fprintln(stderr, name+"を入力してください")
}

func resolveId(id string) (int, error) {
	resolveId, err := strconv.Atoi(id)
	if err != nil {
		return -1, errors.New("idは数値を入力してください: " + id)
	}
	return resolveId, nil
}

func isEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
