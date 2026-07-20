package memo

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func newNotFoundError(id int) error {
	return errors.New("memo not found: " + strconv.Itoa(id))
}

func CreateMemo(memos []Memo, title string, now time.Time) Memo {
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

func FindById(memos []Memo, id int) (Memo, error) {

	for _, memo := range memos {
		if memo.ID == id {
			return memo, nil
		}
	}

	return Memo{}, newNotFoundError(id)
}

func FilterMemos(memos []Memo, keyword string) ([]Memo, error) {
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

func BuildDeletedMemos(memos []Memo, id int) ([]Memo, error) {

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
