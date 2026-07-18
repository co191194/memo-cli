package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateMemo(t *testing.T) {
	now := time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC)

	memos := []Memo{
		{ID: 4, Title: "Exist1"},
		{ID: 1, Title: "Exist2"},
	}

	actual := createMemo(memos, "New Memo", now)

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

	actual, err := findById(memos, 2)
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

	_, actual := findById(memos, 2)
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

	actual, err := filterMemos(memos, "Target")
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

	_, actual := filterMemos(memos, "Target")
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

	actual, err := buildDeletedMemos(memos, 3)
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

	_, actual := buildDeletedMemos(memos, 3)
	if actual == nil {
		t.Fatalf("buildDeletedMemos() any deleted memo")
	}

	assertEqualsMessage(t, actual.Error(), "memo not found: 3")
}

type fakeTimeProvider struct {
}

func (ftp *fakeTimeProvider) Now() time.Time {
	return time.Date(2026, 7, 14, 10, 0, 0, 0, time.Local)
}

func TestAddMemo(t *testing.T) {
	t.Run("バリデーションエラー", func(t *testing.T) {

		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "パラメーターが0個の場合",
				args:     []string{},
				expected: "" +
					"Usage:\n" +
					"  memo add <title>\n",
			},
			{
				testName: "パラメーターが2個の場合",
				args:     []string{"aaa", "bbb"},
				expected: "" +
					"Usage:\n" +
					"  memo add <title>\n",
			},
			{
				testName: "パラメーターが1個で空文字の場合",
				args:     []string{""},
				expected: "タイトルを入力してください\n",
			},
			{
				testName: "パラメーターが1個で空白のみの場合",
				args:     []string{"  "},
				expected: "タイトルを入力してください\n",
			},
		}

		cmd := MemoCommandImpl{}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assertEqualsExitCode(t, cmd.AddMemo(&stdout, &stderr, tc.args), 1)
				assertEqualsMessage(t, stderr.String(), tc.expected)
			})
		}
	})

	t.Run("メモの保存に失敗する場合", func(t *testing.T) {

		cmd := MemoCommandImpl{
			MemoPath:     t.TempDir(),
			TimeProvider: &RealTimeProvider{},
			MemoOperator: &fakeFailSaveMemoOperator{},
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		assertEqualsExitCode(t, cmd.AddMemo(&stdout, &stderr, []string{"Failed Save Memo"}), 1)
		assertEqualsMessage(t, stderr.String(), "メモの保存に失敗しました fake error\n")
	})
}

type fakeFailSaveMemoOperator struct{}

func (f *fakeFailSaveMemoOperator) LoadMemos(path string) ([]Memo, error) {
	return []Memo{}, nil
}

func (f *fakeFailSaveMemoOperator) SaveMemos(path string, memos []Memo) error {
	return errors.New("fake error")
}

func TestListMemos(t *testing.T) {
	t.Run("バリデーションエラー", func(t *testing.T) {

		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "パラメーターが1個の場合",
				args:     []string{"aaaaa"},
				expected: "" +
					"Usage:\n" +
					"  memo list\n",
			},
		}

		cmd := MemoCommandImpl{}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assertEqualsExitCode(t, cmd.ListMemos(&stdout, &stderr, tc.args), 1)
				assertEqualsMessage(t, stderr.String(), tc.expected)
			})
		}
	})
}

var realMemoOperator = MemoOperatorImpl{}

func TestAddMemoAndListMemo(t *testing.T) {

	t.Run("メモが空の場合", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "memos.json")

		var stdout1 bytes.Buffer
		var stderr1 bytes.Buffer

		cmd := MemoCommandImpl{
			MemoPath:     filePath,
			TimeProvider: &fakeTimeProvider{},
			MemoOperator: &realMemoOperator,
		}

		assertEqualsExitCode(t, cmd.ListMemos(&stdout1, &stderr1, []string{}), 0)
		assertEqualsMessage(t, stdout1.String(), "No memos found.\n")

		var stdout2 bytes.Buffer
		var stderr2 bytes.Buffer
		assertEqualsExitCode(t, cmd.AddMemo(&stdout2, &stderr2, []string{"Add First"}), 0)

		var stdout3 bytes.Buffer
		var stderr3 bytes.Buffer
		assertEqualsExitCode(t, cmd.ListMemos(&stdout3, &stderr3, []string{}), 0)
		assertEqualsMessage(t, stdout3.String(), "1 Add First\t2026-07-14\n")
	})

	t.Run("既存のメモがある場合", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "memos.json")

		cmd := MemoCommandImpl{
			MemoPath:     filePath,
			TimeProvider: &fakeTimeProvider{},
			MemoOperator: &realMemoOperator,
		}

		existTime := time.Date(2026, 7, 1, 9, 30, 0, 0, time.Local)
		memos := []Memo{
			{ID: 1, Title: "Exist Memo", Body: "", CreatedAt: existTime, UpdatedAt: existTime},
		}

		if err := realMemoOperator.SaveMemos(filePath, memos); err != nil {
			t.Fatalf("SaveMemos() err = %v", err)
		}

		var stdout1 bytes.Buffer
		var stderr1 bytes.Buffer

		assertEqualsExitCode(t, cmd.AddMemo(&stdout1, &stderr1, []string{"Add Second"}), 0)

		var stdout2 bytes.Buffer
		var stderr2 bytes.Buffer

		assertEqualsExitCode(t, cmd.ListMemos(&stdout2, &stderr2, []string{}), 0)

		assertEqualsMessage(
			t,
			stdout2.String(),
			""+
				"1 Exist Memo\t2026-07-01\n"+
				"2 Add Second\t2026-07-14\n",
		)

	})

}

func TestShowMemo(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "memos.json")
	memos := []Memo{
		{ID: 1, Title: "Other"},
		{
			ID:        2,
			Title:     "Target",
			Body:      "Test Body",
			CreatedAt: time.Date(2026, 7, 15, 10, 15, 30, 45, time.Local),
			UpdatedAt: time.Date(2026, 8, 20, 11, 12, 13, 14, time.Local),
		},
		{ID: 3, Title: "Other"},
	}

	if err := realMemoOperator.SaveMemos(filePath, memos); err != nil {
		t.Fatalf("SaveMemos() err = %v", err)
	}

	cmd := MemoCommandImpl{
		MemoPath:     filePath,
		TimeProvider: &fakeTimeProvider{},
		MemoOperator: &realMemoOperator,
	}

	t.Run("指定IDのメモが存在する場合", func(t *testing.T) {

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		exitCode := cmd.ShowMemo(&stdout, &stderr, []string{"2"})
		assertEqualsExitCode(t, exitCode, 0)

		expected := "" +
			"# Target\n" +
			"\n" +
			"ID: 2\n" +
			"Created: 2026-07-15 10:15\n" +
			"Updated: 2026-08-20 11:12\n" +
			"\n" +
			"Test Body\n"

		assertEqualsMessage(t, stdout.String(), expected)
	})

	t.Run("指定IDのメモが存在しない場合", func(t *testing.T) {

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		exitCode := cmd.ShowMemo(&stdout, &stderr, []string{"4"})

		assertEqualsExitCode(t, exitCode, 1)
		assertEqualsMessage(t, stderr.String(), "memo not found: 4\n")
	})

	t.Run("バリデーションエラー", func(t *testing.T) {

		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "パラメーターが0個の場合",
				args:     []string{},
				expected: "" +
					"Usage:\n" +
					"  memo show <id>\n",
			},
			{
				testName: "パラメーターが2個の場合",
				args:     []string{"aaa", "bbb"},
				expected: "" +
					"Usage:\n" +
					"  memo show <id>\n",
			},
			{
				testName: "パラメーターが数値でない場合",
				args:     []string{"a"},
				expected: "idは数値を入力してください: a\n",
			},
		}

		cmd := MemoCommandImpl{}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assertEqualsExitCode(t, cmd.ShowMemo(&stdout, &stderr, tc.args), 1)
				assertEqualsMessage(t, stderr.String(), tc.expected)
			})
		}
	})

}

func TestSearch(t *testing.T) {

	cmd := MemoCommandImpl{
		MemoPath:     filepath.Join(t.TempDir(), "memos.json"),
		TimeProvider: &fakeTimeProvider{},
		MemoOperator: &realMemoOperator,
	}

	memos := []Memo{
		{
			ID:        1,
			Title:     "Jack",
			Body:      "Blue",
			CreatedAt: time.Date(2026, 1, 1, 9, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2026, 1, 2, 10, 15, 30, 45, time.Local),
		},
		{
			ID:        2,
			Title:     "Nick",
			Body:      "Red",
			CreatedAt: time.Date(2026, 2, 1, 9, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2026, 2, 2, 10, 15, 30, 45, time.Local),
		},
		{
			ID:        3,
			Title:     "Red",
			Body:      "Green",
			CreatedAt: time.Date(2026, 3, 1, 9, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2026, 3, 2, 10, 15, 30, 45, time.Local),
		},
	}

	if err := realMemoOperator.SaveMemos(cmd.MemoPath, memos); err != nil {
		t.Fatalf("SaveMemos() err = %v", err)
	}

	t.Run("バリデーションエラー", func(t *testing.T) {
		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "パラメーターが0個の場合",
				args:     []string{},
				expected: "" +
					"Usage:\n" +
					"  memo search <keyword>\n",
			},
			{
				testName: "パラメーターが2個の場合",
				args:     []string{"aaa", "bbb"},
				expected: "" +
					"Usage:\n" +
					"  memo search <keyword>\n",
			},
			{
				testName: "キーワードが空文字の場合",
				args:     []string{""},
				expected: "キーワードを入力してください\n",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assertEqualsExitCode(t, cmd.SearchMemos(&stdout, &stderr, tc.args), 1)
				assertEqualsMessage(t, stderr.String(), tc.expected)
			})
		}
	})

	t.Run("キーワードに該当するメモが存在する場合", func(t *testing.T) {

		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "keyword = Red",
				args:     []string{"Red"},
				expected: "" +
					"2 Nick\t2026-02-01\n" +
					"3 Red\t2026-03-01\n",
			},
			{
				testName: "keyword = Jack",
				args:     []string{"Jack"},
				expected: "" +
					"1 Jack\t2026-01-01\n",
			},
			{
				testName: "keyword = e",
				args:     []string{"e"},
				expected: "" +
					"1 Jack\t2026-01-01\n" +
					"2 Nick\t2026-02-01\n" +
					"3 Red\t2026-03-01\n",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assertEqualsExitCode(t, cmd.SearchMemos(&stdout, &stderr, tc.args), 0)
				assertEqualsMessage(t, stdout.String(), tc.expected)
			})
		}
	})

	t.Run("キーワードに該当するメモが存在しない場合", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		assertEqualsExitCode(t, cmd.SearchMemos(&stdout, &stderr, []string{"Yellow"}), 1)
		assertEqualsMessage(t, stderr.String(), "No matching memos found.\n")
	})
}

func TestDeleteMemo(t *testing.T) {

	t.Run("指定IDのメモが存在する場合", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "memos.json")

		memos := []Memo{
			{ID: 1, Title: "Other"},
			{ID: 2, Title: "Target"},
			{ID: 3, Title: "Other"},
		}

		if err := realMemoOperator.SaveMemos(filePath, memos); err != nil {
			t.Fatalf("SaveMemos() err = %v", err)
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd := MemoCommandImpl{
			MemoPath:     filePath,
			TimeProvider: &fakeTimeProvider{},
			MemoOperator: &realMemoOperator,
		}

		assertEqualsExitCode(t, cmd.DeleteMemo(&stdout, &stderr, []string{"2"}), 0)

		actual, err := realMemoOperator.LoadMemos(filePath)
		if err != nil {
			t.Fatalf("LoadMemos() err = %v", err)
		}

		if len(actual) != 2 {
			t.Fatalf("actual.len = %d, expected = 2", len(actual))
		}

		if actual[0].ID != 1 {
			t.Errorf("actual[0].ID = %d, expected = 1", actual[0].ID)
		}

		if actual[1].ID != 3 {
			t.Errorf("actual[1].ID = %d, expected = 3", actual[1].ID)
		}
	})

	t.Run("指定IDのメモが存在しない場合", func(t *testing.T) {

		filePath := filepath.Join(t.TempDir(), "memos.json")

		memos := []Memo{
			{ID: 1, Title: "Other"},
			{ID: 3, Title: "Other"},
		}

		if err := realMemoOperator.SaveMemos(filePath, memos); err != nil {
			t.Fatalf("SaveMemos() err = %v", err)
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd := MemoCommandImpl{
			MemoPath:     filePath,
			TimeProvider: &fakeTimeProvider{},
			MemoOperator: &realMemoOperator,
		}

		exitCode := cmd.DeleteMemo(&stdout, &stderr, []string{"2"})
		assertEqualsExitCode(t, exitCode, 1)

		expected := "memo not found: 2\n"
		assertEqualsMessage(t, stderr.String(), expected)

		actual, err := realMemoOperator.LoadMemos(filePath)
		if err != nil {
			t.Fatalf("LoadMemos() err = %v", err)
		}

		if len(actual) != 2 {
			t.Fatalf("actual.len = %d, expected = 2", len(actual))
		}

	})

	t.Run("バリデーションエラー", func(t *testing.T) {

		testCases := []struct {
			testName string
			args     []string
			expected string
		}{
			{
				testName: "パラメーターが0個の場合",
				args:     []string{},
				expected: "" +
					"Usage:\n" +
					"  memo delete <id>\n",
			},
			{
				testName: "パラメーターが2個の場合",
				args:     []string{"1", "aaa"},
				expected: "" +
					"Usage:\n" +
					"  memo delete <id>\n",
			},
			{
				testName: "パラメーターが数値でない場合",
				args:     []string{"a"},
				expected: "idは数値を入力してください: a\n",
			},
		}

		cmd := MemoCommandImpl{
			MemoPath:     "",
			TimeProvider: &fakeTimeProvider{},
		}

		for _, tc := range testCases {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := cmd.DeleteMemo(&stdout, &stderr, tc.args)

			assertEqualsExitCode(t, exitCode, 1)
			assertEqualsMessage(t, stderr.String(), tc.expected)
		}
	})

}

func assertEqualsExitCode(t *testing.T, actual int, expected int) {
	t.Helper()
	if actual != expected {
		t.Fatalf("actual ExitCode = %d, expected = %d", actual, expected)
	}
}

func assertEqualsMessage(t *testing.T, actual string, expected string) {
	t.Helper()
	if actual != expected {
		t.Errorf("actual = %q, expected = %q", actual, expected)
	}
}
