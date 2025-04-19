# puki2md

puki2mdはPukiWikiの記事をMarkdown形式に変換するパッケージです。

## インストール

```
go get github.com/miso24/puki2md
```

## 使用方法

```go
package main

import (
    "fmt"
    "strings"

    "github.com/miso24/puki2md"
) 

func main() {
    article := `* 見出し1

''強調''や[[リンク]]に対応しています。

- リスト

| a | b |
| c | d |
`
    
    md := puki2md.Convert(strings.NewReader(article))
    fmt.Println(md)
}
```

## 対応状況

### ブロック要素

|PukiWiki記法|対応状況|
|---|---|
|見出し `* Heading`|✔|
|引用 `> quote`|❌️|
|リスト `- list`|✔|
|番号付きリスト `+ list`|✔|
|定義リスト `: term \| desc`|✔|
|表組み `\|data\|data\|`|✔|
|csv形式の表組み `,data,data,...`|✔|
|添付ファイル `#ref(...)`|❌️|

### インライン要素

|PukiWiki記法|対応状況|
|---|---|
|強調 `''text''`|✔|
|斜体 `'''text'''`|✔|
|サイズ `&size(...){...};`|❌️|
|文字色 `&color(...){...};`|❌️|
|取り消し線 `%%text%%`|❌️|
|注釈 `((text))`|✔|
|添付ファイル `&ref(...);`|🚧|
|リンク `[[link]]`|✔|

> ✔: 実装済み  ❌️: 未対応  🚧:開発中
