# Golang batch library

指定した時間に指定したコマンドを実行するやつ．

グループごとに同時実行数を管理できます．

とりあえず動かすためのやっつけ実装です．．．まともな実装になることは永久に無い気がするのでプルリクエスト待ってます．

TODO: RPCとかデーモン動作とかコマンドラインツールとか．

## Installation

``` bash
go get "github.com/binzume/gobatch"

```

## Usage


``` go
package main

import "time"
import "github.com/binzume/gobatch"

func main() {
	batch := gobatch.Default;
	batch.RegisterGroup("hoge", 2, gobatch.RunAt)

	batch.CommandAt("hoge", "job00", "sleep 100", time.Now())
	batch.CommandAt("hoge", "job01", "sleep 13", time.Now().Add(10*time.Second))
	batch.CommandAt("hoge", "job02", "sleep 12", time.Now().Add(4*time.Second))
	batch.CommandAt("hoge", "job03", "sleep 11", time.Now().Add(4*time.Second))
	batch.CommandAt("hoge", "job04", "sleep 11", time.Now().Add(4*time.Second))
	batch.CancelById("hoge", "job03");

	time.Sleep(20*time.Second)
	batch.CancelById("hoge", "job00")

	batch.Wait()
}
```


- batch := gobatch.Default // デフォルトのBatchインスタンス
- batch.RegisterGroup("hoge", 2, gobatch.RunAt) // バッチグループの登録．同時実行制限2，指定時間になったら実行
- batch.CommandAt("hoge", "job01", "sleep 13", time.Now().Add(10*time.Second)) // 10秒後にsleepコマンドを実行
- batch.CancelById("hoge", "job00") // バッチをキャンセル．実行されていたらKillする．
- batch.Wait() // バッチインスタンスの終了を待つ (現状は無限に待ちます)


## License:

- MIT License

