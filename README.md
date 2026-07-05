# memo

Go学習用のシンプルなMarkdownメモCLIです。

ターミナルからメモを追加・一覧表示・詳細表示・検索・削除できる個人用ツールを作成します。

このプロジェクトは、Goの基本文法や標準ライブラリの使い方を学ぶことを目的としています。

## 目的

このアプリケーションの主な目的は、Goで小さなCLIアプリを作りながら、以下の内容を学習することです。

* Goの基本文法
* `struct` の定義と利用
* スライスの操作
* ファイルの読み書き
* JSONのエンコード・デコード
* エラーハンドリング
* コマンドライン引数の扱い
* 関数分割
* パッケージ分割
* ユニットテスト

## 概要

`memo` は、ローカル環境でMarkdown形式のメモを管理するためのCLIツールです。

最初のバージョンでは、データベースは使用せず、メモをJSONファイルに保存します。

保存先の例:

```txt
~/.memo/memos.json
```

## 完成イメージ

以下のようなコマンドでメモを管理できることを目指します。

```bash
memo add "Goのstructを学ぶ"
memo add "JSONファイル保存を実装する"
memo list
memo show 1
memo search JSON
memo delete 1
memo list
```

## MVPで実装する機能

最初のバージョンでは、以下の5つの機能を実装します。

| コマンド                    | 説明               |
| ----------------------- | ---------------- |
| `memo add <title>`      | メモを追加する          |
| `memo list`             | メモを一覧表示する        |
| `memo show <id>`        | 指定したIDのメモを詳細表示する |
| `memo search <keyword>` | キーワードでメモを検索する    |
| `memo delete <id>`      | 指定したIDのメモを削除する   |

## コマンド仕様

### メモを追加する

```bash
memo add "Goのエラーハンドリングについて学ぶ"
```

最初の実装では、タイトルのみを指定します。

本文は空文字として保存します。

### メモの一覧を表示する

```bash
memo list
```

出力例:

```txt
1  Goのエラーハンドリングについて学ぶ  2026-07-05
2  JSONファイル保存を実装する          2026-07-05
```

一覧では、以下の情報を表示します。

* ID
* タイトル
* 作成日

### メモの詳細を表示する

```bash
memo show 1
```

出力例:

```txt
# Goのエラーハンドリングについて学ぶ

ID: 1
Created: 2026-07-05 18:00
Updated: 2026-07-05 18:00

本文はここに表示される
```

### メモを検索する

```bash
memo search Go
```

タイトルまたは本文にキーワードが含まれるメモを検索します。

最初の実装では、大文字・小文字を区別して検索します。

### メモを削除する

```bash
memo delete 1
```

指定したIDのメモを削除します。

最初の実装では、削除前の確認メッセージは表示しません。

## データ構造

メモは以下のような構造体で表現します。

```go
type Memo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 保存形式

メモはJSONファイルに保存します。

保存データの例:

```json
[
  {
    "id": 1,
    "title": "Goのエラーハンドリングについて学ぶ",
    "body": "",
    "created_at": "2026-07-05T18:00:00+09:00",
    "updated_at": "2026-07-05T18:00:00+09:00"
  },
  {
    "id": 2,
    "title": "JSONファイル保存を実装する",
    "body": "",
    "created_at": "2026-07-05T18:10:00+09:00",
    "updated_at": "2026-07-05T18:10:00+09:00"
  }
]
```

## セットアップ

### 前提

Goがインストールされていること。

```bash
go version
```

### リポジトリをクローンする

```bash
git clone https://github.com/<your-name>/memo.git
cd memo
```

### モジュールを初期化する

新規作成時は以下を実行します。

```bash
go mod init github.com/<your-name>/memo
```

既に `go.mod` が存在する場合、この手順は不要です。

### ビルドする

```bash
go build -o memo
```

### 実行する

```bash
./memo list
```

## 開発手順

このプロジェクトは、以下の順番で実装します。

### Step 1: プロジェクトを作成する

```bash
mkdir memo
cd memo
go mod init github.com/<your-name>/memo
```

最初のファイル構成:

```txt
memo/
  go.mod
  main.go
  memo.go
  store.go
  memo_test.go
  store_test.go
```

### Step 2: `Memo` 構造体を定義する

まず、1件のメモを表す `Memo` 構造体を作成します。

```go
type Memo struct {
    ID        int
    Title     string
    Body      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Step 3: メモをメモリ上で扱う

最初はファイル保存を考えず、スライスで複数のメモを扱います。

実装する処理:

* メモを追加する
* IDでメモを探す
* キーワードでメモを検索する
* IDでメモを削除する

### Step 4: JSONファイルに保存する

メモをJSONファイルへ保存・読み込みできるようにします。

作成する関数の例:

```go
func LoadMemos(path string) ([]Memo, error)
func SaveMemos(path string, memos []Memo) error
```

### Step 5: `add` コマンドを実装する

```bash
memo add "最初のメモ"
```

処理の流れ:

1. JSONファイルから既存のメモを読み込む
2. 新しいIDを採番する
3. 新しいメモを作成する
4. メモ一覧に追加する
5. JSONファイルへ保存する

IDの採番は、最初は「現在の最大ID + 1」とします。

### Step 6: `list` コマンドを実装する

```bash
memo list
```

処理の流れ:

1. JSONファイルからメモを読み込む
2. メモの一覧を表示する

### Step 7: `show` コマンドを実装する

```bash
memo show 1
```

指定されたIDのメモを検索し、詳細を表示します。

対象のメモが存在しない場合は、エラーメッセージを表示します。

```txt
memo not found: 1
```

### Step 8: `search` コマンドを実装する

```bash
memo search Go
```

タイトルまたは本文にキーワードが含まれるメモを表示します。

最初は `strings.Contains` を使って実装します。

### Step 9: `delete` コマンドを実装する

```bash
memo delete 1
```

指定されたID以外のメモだけを残し、JSONファイルへ保存し直します。

対象のメモが存在しない場合は、エラーとして扱います。

### Step 10: テストを書く

Go標準のテスト機能を使ってユニットテストを書きます。

```bash
go test ./...
```

テスト対象の例:

* メモ追加時にIDが採番されること
* IDでメモを取得できること
* キーワードで検索できること
* メモを削除できること
* JSONファイルに保存できること
* JSONファイルから読み込めること

## ディレクトリ構成

最初はシンプルな構成で実装します。

```txt
memo/
  go.mod
  main.go
  memo.go
  store.go
  memo_test.go
  store_test.go
```

コードが増えてきたら、以下のような構成にリファクタリングします。

```txt
memo/
  go.mod
  cmd/
    memo/
      main.go
  internal/
    memo/
      memo.go
      service.go
    store/
      json_store.go
```

## 最初は実装しない機能

以下の機能は、MVP完成後に必要に応じて追加します。

* MarkdownのHTMLプレビュー
* タグ機能
* カテゴリ機能
* エディタ連携
* 複数ファイル管理
* SQLite対応
* Git連携
* Web UI
* TUI表示
* 暗号化

## 今後の拡張案

MVP完成後は、以下の機能を追加していく予定です。

| 機能             | 学べること          |
| -------------- | -------------- |
| `edit` コマンド    | 既存データの更新処理     |
| `--body` オプション | コマンドラインオプション処理 |
| タグ機能           | スライス、検索条件      |
| 設定ファイル         | OSごとのパス管理      |
| SQLite対応       | データベースアクセス     |
| Markdownファイル出力 | ファイル生成         |
| テスト強化          | 設計改善           |
| GitHub Actions | CI             |

## 発展版のコマンド例

将来的には、以下のようなコマンドに拡張できます。

```bash
memo add "Goのメモ" --body "Goではエラーを戻り値として扱う"
memo edit 1
memo list --limit 10
memo search error --case-insensitive
memo export 1 --format markdown
memo open 1
```

## 完成条件

MVPの完成条件は、以下のコマンドがすべて動作することです。

```bash
memo add "Goのstructを学ぶ"
memo add "JSONファイル保存を実装する"
memo list
memo show 1
memo search JSON
memo delete 1
memo list
```

加えて、メモがJSONファイルに永続化されていることを確認します。

## ライセンス

MIT License
