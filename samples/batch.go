package main

import "time"
import "github.com/binzume/gobatch"

func main() {
	batch := gobatch.Default;
	batch.RegisterGroup("hoge", 2)

	batch.CommandAt("hoge", "sleep 13", time.Now().Add(10*time.Second))
	batch.CommandAt("hoge", "sleep 12", time.Now().Add(3*time.Second))
	batch.CommandAt("hoge", "sleep 11", time.Now().Add(4*time.Second))

	batch.Wait()
}

