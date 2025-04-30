package models

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleIsEmpty
	}

	if r.Description == "" {
		return ErrDescriptionIsEmpty
	}

	if r.Status == "" {
		return ErrStatusIsEmpty
	}

	return nil
}

func (r *CreateTaskRequest) ConvertToTask() *Task {
	return &Task{
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
	}
}

func (r *UpdateTaskRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleIsEmpty
	}

	if r.Description == "" {
		return ErrDescriptionIsEmpty
	}

	if r.Status == "" {
		return ErrStatusIsEmpty
	}

	return nil
}

func (r *UpdateTaskRequest) ConvertToTask(id string) *Task {
	return &Task{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
	}
}
