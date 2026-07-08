package main

import (
	"fmt"
	"os"
	"time"
)

const path = "~/.memo/memos.json"

func AddMemo(title string) {

	memos, err := LoadMemos(path)
	if err != nil {
		fmt.Println("メモを開くことができませんでした", err)
		os.Exit(1)
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
