package todo

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddTask(task *Task) error {
	store, err := s.storage.Load()
	if err != nil {
		return err
	}

	store.Tasks = append(store.Tasks, *task)
	return s.storage.Save(store)
}

func (s *Service) GetAllTasks() ([]Task, error) {
	store, err := s.storage.Load()
	if err != nil {
		return nil, err
	}
	return store.Tasks, nil
}

func (s *Service) GetTaskByID(id string) (*Task, error) {
	store, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	for i, task := range store.Tasks {
		if task.ID == id {
			return &store.Tasks[i], nil
		}
	}
	return nil, nil
}

func (s *Service) UpdateTask(id string, task *Task) error {
	store, err := s.storage.Load()
	if err != nil {
		return err
	}

	for i, t := range store.Tasks {
		if t.ID == id {
			store.Tasks[i] = *task
			return s.storage.Save(store)
		}
	}
	return nil
}

func (s *Service) DeleteTask(id string) error {
	store, err := s.storage.Load()
	if err != nil {
		return err
	}

	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks = append(store.Tasks[:i], store.Tasks[i+1:]...)
			return s.storage.Save(store)
		}
	}
	return nil
}

func (s *Service) FilterTasks(status Status, priority Priority, tag string) ([]Task, error) {
	tasks, err := s.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var filtered []Task
	for _, task := range tasks {
		if status != "" && task.Status != status {
			continue
		}
		if priority != "" && task.Priority != priority {
			continue
		}
		if tag != "" {
			found := false
			for _, t := range task.Tags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filtered = append(filtered, task)
	}

	return filtered, nil
}
