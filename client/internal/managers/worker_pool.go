package managers

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type Task struct {
	Ctx     context.Context
	Command string
	Args    []string
}
type Worker struct {
	id   int
	quit chan bool
}

type WorkerPool struct {
	TaskChan chan Task
	Workers  []Worker
	wg       *sync.WaitGroup
	cm       *CommandManager
	mu       sync.Mutex
}

func NewWorkerPool(taskChan chan Task, poolSize int, cm *CommandManager) *WorkerPool {
	workerPool := &WorkerPool{
		TaskChan: taskChan,
		Workers:  make([]Worker, 0),
		wg:       &sync.WaitGroup{},
		cm:       cm,
		mu:       sync.Mutex{},
	}
	workerPool.SetWorkers(poolSize)
	return workerPool
}

func (w *WorkerPool) SetWorkers(numberOfWorkers int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	currentWorkers := len(w.Workers)
	if numberOfWorkers >= currentWorkers {
		for i := currentWorkers; i < numberOfWorkers; i++ {
			worker := Worker{id: i + 1, quit: make(chan bool)}
			w.Workers = append(w.Workers, worker)
			w.wg.Add(1)
			go w.worker(worker)
		}
		fmt.Printf("Added %d workers, total: %d\n", len(w.Workers)-currentWorkers, len(w.Workers))
	} else if numberOfWorkers < currentWorkers {
		// в данной ситуации я сомневаюсь что это хорошая идея, но так, я увижу сообщение о закрытии от
		// той горутины, которая обрабатывала эту команду
		go func() {
			for i := currentWorkers - 1; i > numberOfWorkers; i-- {
				w.Workers[i].quit <- true
				w.Workers = w.Workers[:i]
			}
			fmt.Printf("Removed %d workers, total: %d\n", currentWorkers-len(w.Workers), len(w.Workers))
		}()
	}
}
func (w *WorkerPool) worker(worker Worker) {
	defer w.wg.Done()

	for {
		select {
		case task := <-w.TaskChan:
			fmt.Printf("Worker %d processing command: %s\n", worker.id, task.Command)
			if task.Command == "set-workers" {
				numberOfWorkers, err := strconv.Atoi(task.Args[0])
				if err != nil {
					fmt.Printf("Error converting arguments to int64: %s\n", task.Args[0])
				}
				w.SetWorkers(numberOfWorkers)
			} else if err := w.cm.ExecuteCommand(task); err != nil {
				fmt.Println(err)
			}
		case <-worker.quit:
			fmt.Printf("Worker %d quitting\n", worker.id)
			close(worker.quit)
			return
		}
	}
}

func (w *WorkerPool) Close() {
	go func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		for _, worker := range w.Workers {
			worker.quit <- true
		}
	}()
	w.wg.Wait()
	close(w.TaskChan)
}
