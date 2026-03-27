package task

import (
	"errors"
	"sync"

	"aiguard/internal/uiapi"
)

type TaskHandle struct {
	Cancel func()
	Req    uiapi.StartReviewRequest
	State  TaskState
}

type Manager struct {
	mu    sync.Mutex
	tasks map[string]*TaskHandle
}

type TaskState string

const (
	TaskStateRunning    TaskState = "running"
	TaskStateCancelling TaskState = "cancelling"
)

func NewManager() *Manager {
	return &Manager{
		tasks: map[string]*TaskHandle{},
	}
}

func (m *Manager) Add(taskID string, cancel func(), req uiapi.StartReviewRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[taskID] = &TaskHandle{
		Cancel: cancel,
		Req:    req,
		State:  TaskStateRunning,
	}
}

func (m *Manager) Done(taskID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tasks, taskID)
}

func (m *Manager) Cancel(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	handle, ok := m.tasks[taskID]
	if !ok {
		return errors.New("任务不存在或已结束")
	}
	if handle.State == TaskStateCancelling {
		return nil
	}

	handle.State = TaskStateCancelling
	handle.Cancel()
	return nil
}

func (m *Manager) HasRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.tasks) > 0
}
