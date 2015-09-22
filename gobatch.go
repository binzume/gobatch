package gobatch

import (
	"container/heap"
    "container/list"
	"os/exec"
	"time"
	"log"
)

type TaskFunc func() (result int, quit <-chan int)

type Task struct {
	Id string
	Group string
	Type int
	priority int64
	index int
	result int
	v string
	C <-chan *Task
	f TaskFunc
}

type BatchQueue []*Task

func (pq BatchQueue) Len() int { return len(pq) }

func (pq BatchQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq BatchQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *BatchQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Task)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *BatchQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *BatchQueue) Last() *Task {
	n := len(*pq)
	return (*pq)[n-1]
}

// update modifies the priority and value of an Item in the queue.
func (pq *BatchQueue) update(item *Task, value string, priority int64) {
	// item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

type BatchGroup struct{
	Name string
	NumConcurrent int
	Type int
	running *list.List // list of *Task
	ticker *time.Ticker
}

type Batch struct{
	queue map[string] *BatchQueue
	group map[string] *BatchGroup
	exit chan int
}

func (b *Batch) AddGroup(g *BatchGroup) {
	b.queue[g.Name] = &BatchQueue{}
	b.group[g.Name] = g
	heap.Init(b.queue[g.Name])

	// check each seconds.
	g.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
				case <- g.ticker.C:
					// log.Printf(g.Name + " running %d tasks.", g.running.Len())
					for b.queue[g.Name].Len() > 0 && g.running.Len() < g.NumConcurrent {
						task := heap.Pop(b.queue[g.Name]).(*Task)
						if task.priority <= time.Now().Unix() {
							log.Print("start " + task.Id )
							if task.f != nil {
								_, done := task.f()
								e := g.running.PushBack(task)
								go func(){
									<- done
									g.running.Remove(e)
									//for t := g.running.Front(); t != nil; t = t.Next() {
									//	if t.Value.(*Task) == task {
									//		break
									//	}
									//}
								}()
							}
						} else {
							// requeue
							heap.Push(b.queue[g.Name], task)
							break
						}
					}
			}
		}
	}()
}

func (b *Batch) RegisterGroup(name string, numConcurrent int) {
	b.AddGroup(&BatchGroup{Name: name, NumConcurrent: numConcurrent, Type: 0, running: list.New()});
}

func (b *Batch) AddTask(t *Task) {
	heap.Push(b.queue[t.Group], t)
}

func (b *Batch) GetRunnings(group string) []*Task {
	ts := []*Task{}
	for t := b.group[group].running.Front(); t != nil; t = t.Next() {
		ts = append(ts, t.Value.(*Task))
	}
	return ts
}

func New() *Batch {
	b := &Batch{
		queue: map[string] *BatchQueue{},
		group: map[string] *BatchGroup{},
		exit: make(chan int, 1)}
	log.Print("Start batch.")
	return b;
}

var Default *Batch = New()

func (b *Batch) CommandAt(group string, cmd string, at time.Time) *Task {
	task := &Task{Id: "", v: cmd, Group: group, priority: at.Unix()}
	b.AddTask(task)
	task.f = func() (r int, q <- chan int) {
		log.Print("exec " + cmd)
		done := make(chan int, 1)
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Start()
		go func(){
			cmd.Wait()
			done <- 1
		}()
		return 0, done
	}
	return task
}


func (b *Batch) Wait() {
	<- b.exit
}

