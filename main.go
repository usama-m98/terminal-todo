package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var InvalidArgs error = errors.New(
	fmt.Sprintln(
		"\n\n",
		"Too many arguments provided\n\n",
		"Available commands:\n\n",
		"todo add <task>\n",
		"todo delete <task_id>\n",
		"todo list\n",
		"todo complete <task_id>",
	),
)

type TaskConfig struct {
	DB *TaskStore
}

func NewTaskConfig(ts *TaskStore) TaskConfig {

	return TaskConfig{
		DB: ts,
	}
}

var count int = 0

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var err error
	if len(os.Args) > 4 {
		return InvalidArgs
	}

	args := os.Args[1:]

	taskStr, err := NewTaskStore()
	if err != nil {
		return err
	}
	taskCfg := NewTaskConfig(taskStr)

	switch args[0] {
	case "add":
		err = taskCfg.addTask(args)
	case "list":
		taskCfg.listAll()
	case "complete":
		err = taskCfg.completeTask(args)
	case "delete":
		err = taskCfg.removeTask(args)
	default:
		return InvalidArgs
	}

	if err != nil {
		return err
	}

	return nil
}

func (c TaskConfig) addTask(args []string) error {
	if len(args) == 1 {
		return errors.New("Task is required")
	}

	id, err := c.DB.addToStore(args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Task id: %d", id)
	return nil
}

func (c TaskConfig) completeTask(args []string) error {
	if len(args) == 1 {
		return errors.New("task id is required")
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("Task ids can only be numbers")
	}

	err = c.DB.markComplete(id)
	if err != nil {
		return errors.New("Failed to make complete")
	}

	c.listAll()

	return nil
}

func (c TaskConfig) removeTask(args []string) error {
	if len(args) == 1 {
		return errors.New("task id is required")
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("Task ids can only be numbers")
	}

	err = c.DB.removeFromStore(id)
	if err != nil {
		return errors.New("Failed to delete task from task store")
	}

	fmt.Printf("Removed task id: %d\n", id)

	return nil
}

func (c TaskConfig) listAll() {
	list := "ID\tTask\tStatus\n"

	for k, v := range c.DB.store {
		list += fmt.Sprintf("%d\t%s\t%s\n", k, v.Title, v.Status)
	}

	fmt.Println(list)
}
