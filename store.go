package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const path = "~/.memo/memos.json"

func AddMemo(title string) {

	memos, err := LoadMemos(path)
	if err != nil {
		printOpenFileError(err)
	}

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

	now := time.Now()

	addedMemo := Memo{
		ID:        id,
		Title:     title,
		Body:      "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	memos = append(memos, addedMemo)

	if err := SaveMemos(path, memos); err != nil {
		fmt.Println("メモの保存に失敗しました", err)
		os.Exit(1)
	}
}

func ListMemos() {
	memos, err := LoadMemos(path)
	if err != nil {
		printOpenFileError(err)
	}

	if len(memos) == 0 {
		fmt.Println("No memos found.")
	} else {
		for _, memo := range memos {
			fmt.Printf("%d %s\t%s\n", memo.ID, memo.Title, memo.CreatedAt.Format("2006-01-02"))
		}
	}
}

const DATE_TIME_FORMAT = "2006-01-02 15:04"

func ShowMemo(id int) {
	memos, err := LoadMemos(path)
	if err != nil {
		printOpenFileError(err)
	}

	for _, memo := range memos {
		if memo.ID == id {
			fmt.Println("# " + memo.Title)
			fmt.Println()
			fmt.Println("ID: " + strconv.Itoa(memo.ID))
			fmt.Println("Created: " + memo.CreatedAt.Format(DATE_TIME_FORMAT))
			fmt.Println("Updated: " + memo.UpdatedAt.Format(DATE_TIME_FORMAT))
			fmt.Println()
			fmt.Println(memo.Body)
			return
		}
	}
	fmt.Printf("memo not found: %d\n", id)
	os.Exit(1)
}

func printOpenFileError(err error) {
	fmt.Println("メモを開くことができませんでした", err)
	os.Exit(1)
}
