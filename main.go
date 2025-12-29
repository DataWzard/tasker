package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	DoneAt    time.Time `json:"done_at,omitempty"`
}

const dataFile = "tasks.json"

func loadTasks() ([]Task, error) {
	b, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	if len(b) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, b, 0644)
}

func nextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func usage() {
	fmt.Println("tasker: simple todo CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  tasker add <title>")
	fmt.Println("  tasker list")
	fmt.Println("  tasker done <id>")
	fmt.Println("  tasker rm <id>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		os.Exit(1)
	}

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Missing title.")
			os.Exit(1)
		}
		title := strings.Join(os.Args[2:], " ")
		task := Task{
			ID:        nextID(tasks),
			Title:     title,
			Done:      false,
			CreatedAt: time.Now(),
		}
		tasks = append(tasks, task)
		if err := saveTasks(tasks); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}
		fmt.Printf("Added #%d: %s\n", task.ID, task.Title)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("No tasks yet.")
			return
		}
		for _, t := range tasks {
			box := " "
			if t.Done {
				box = "x"
			}
			fmt.Printf("[%s] #%d %s\n", box, t.ID, t.Title)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Missing id.")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid id.")
			os.Exit(1)
		}
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				tasks[i].DoneAt = time.Now()
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Task not found.")
			os.Exit(1)
		}
		if err := saveTasks(tasks); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}
		fmt.Println("Marked done:", id)

	case "rm":
		if len(os.Args) < 3 {
			fmt.Println("Missing id.")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid id.")
			os.Exit(1)
		}
		out := make([]Task, 0, len(tasks))
		found := false
		for _, t := range tasks {
			if t.ID == id {
				found = true
				continue
			}
			out = append(out, t)
		}
		if !found {
			fmt.Println("Task not found.")
			os.Exit(1)
		}
		if err := saveTasks(out); err != nil {
			fmt.Println("Error saving tasks:", err)
			os.Exit(1)
		}
		fmt.Println("Removed:", id)

	default:
		usage()
		os.Exit(1)
	}
}
