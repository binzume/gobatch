package main

import "time"
import "github.com/binzume/gobatch"

func main() {
	batch := gobatch.Default;
	batch.RegisterGroup("hoge", 2)

	batch.CommandAt("hoge", "job00", "sleep 100", time.Now())
	batch.CommandAt("hoge", "job01", "sleep 13", time.Now().Add(10*time.Second))
	batch.CommandAt("hoge", "job02", "sleep 12", time.Now().Add(4*time.Second))
	batch.CommandAt("hoge", "job03", "sleep 11", time.Now().Add(4*time.Second))
	batch.CommandAt("hoge", "job04", "sleep 11", time.Now().Add(4*time.Second))
	batch.CancelById("hoge", "job03");

	time.Sleep(20*time.Second)
	batch.CancelById("hoge", "job00");

	batch.Wait()
}
