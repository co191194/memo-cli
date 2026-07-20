package memo_test

import (
	"testing"
	"time"

	"github.com/co191194/memo-cli/internal/memo"
)

type Memo = memo.Memo

func TestCreateMemo(t *testing.T) {
	now := time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC)

	memos := []Memo{
		{ID: 4, Title: "Exist1"},
		{ID: 1, Title: "Exist2"},
	}

	actual := memo.CreateMemo(memos, "New Memo", now)

	if actual.ID != 5 {
		t.Errorf("ID = %d, expected = 5", actual.ID)
	}

	if actual.Title != "New Memo" {
		t.Errorf("Title = %q, expected = %q", actual.Title, "New Memo")
	}

	if !actual.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, expected = %v", actual.CreatedAt, now)
	}

	if !actual.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, expected = %v", actual.UpdatedAt, now)
	}
}

func TestFindById(t *testing.T) {
	memos := []Memo{
		{ID: 1, Title: "Other"},
		{ID: 3, Title: "Other"},
		{ID: 2, Title: "Target"},
		{ID: 4, Title: "Other"},
	}

	actual, err := memo.FindById(memos, 2)
	if err != nil {
		t.Fatalf("findById() err = %v", err)
	}

	if actual.ID != 2 {
		t.Errorf("ID = %d, expected = 2", actual.ID)
	}

	if actual.Title != "Target" {
		t.Errorf("Title = %q, expected = %q", actual.Title, "Target")
	}
}

func TestFindById_NotFound(t *testing.T) {
	memos := []Memo{
		{ID: 1, Title: "Other"},
		{ID: 3, Title: "Other"},
	}

	_, actual := memo.FindById(memos, 2)
	if actual == nil {
		t.Fatalf("findById any match")
	}

	assertEqualsMessage(t, actual.Error(), "memo not found: 2")
}

func TestFilterMemos(t *testing.T) {
	memos := []Memo{
		{ID: 1, Title: "Other Title1", Body: "Other Body1"},
		{ID: 2, Title: "Target Title2", Body: "Body2"},
		{ID: 3, Title: "Title3", Body: "Body3 Target"},
		{ID: 4, Title: "Other Title4", Body: "Other Body4"},
		{ID: 5, Title: "Title5 Target Middle", Body: "Body5"},
		{ID: 6, Title: "Other Title6", Body: "Other Body6"},
	}

	actual, err := memo.FilterMemos(memos, "Target")
	if err != nil {
		t.Fatalf("filterMemos() err = %v", err)
	}

	if len(actual) != 3 {
		t.Fatalf("actual len = %d, expected = 3", len(actual))
	}
	if actual[0].ID != 2 {
		t.Errorf("actual[0].ID = %d, expected = 2", actual[0].ID)
	}
	if actual[1].ID != 3 {
		t.Errorf("actual[1].ID = %d, expected = 3", actual[1].ID)
	}
	if actual[2].ID != 5 {
		t.Errorf("actual[2].ID = %d, expected = 5", actual[2].ID)
	}
}

func TestFilterMemos_NotMatch(t *testing.T) {
	memos := []Memo{
		{ID: 1, Title: "Other Title1", Body: "Other Body1"},
		{ID: 6, Title: "Other Title6", Body: "Other Body6"},
	}

	_, actual := memo.FilterMemos(memos, "Target")
	if actual == nil {
		t.Fatalf("filterMemos any match")
	}

	assertEqualsMessage(t, actual.Error(), "No matching memos found.")
}

func TestBuildDeletedMemos(t *testing.T) {
	memos := []Memo{
		{ID: 6, Title: "Other1"},
		{ID: 3, Title: "Target"},
		{ID: 1, Title: "Other2"},
	}

	actual, err := memo.BuildDeletedMemos(memos, 3)
	if err != nil {
		t.Fatalf("buildDeletedMemos() err = %v", err)
	}

	if len(actual) != 2 {
		t.Fatalf("actual.len = %d, expected = 2", len(actual))
	}

	if actual[0].ID != 6 {
		t.Errorf("actual[0].ID = %d, expected = 6", actual[0].ID)
	}

	if actual[1].ID != 1 {
		t.Errorf("actual[1].ID = %d, expected = 1", actual[1].ID)
	}
}

func TestBuildDeletedMemos_NotFountMemo(t *testing.T) {
	memos := []Memo{
		{ID: 6, Title: "Other1"},
		{ID: 1, Title: "Other2"},
	}

	_, actual := memo.BuildDeletedMemos(memos, 3)
	if actual == nil {
		t.Fatalf("buildDeletedMemos() any deleted memo")
	}

	assertEqualsMessage(t, actual.Error(), "memo not found: 3")
}

func assertEqualsMessage(t *testing.T, actual string, expected string) {
	t.Helper()
	if actual != expected {
		t.Errorf("actual = %q, expected = %q", actual, expected)
	}
}
