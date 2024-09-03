// File: graph/schema.resolvers.go

package graph

import (
	"context"
	"fmt"
	"log"
	"todo_app/db"
	"todo_app/graph/model"
)

func (r *mutationResolver) AddTask(ctx context.Context, title string, description *string, dueDate *string) (*model.Task, error) {
    var dueDateValue string
    if dueDate != nil {
        dueDateValue = *dueDate
    } else {
        dueDateValue = "No due date" // Handle the default case, or set to a default date string
    }

    log.Printf("Received AddTask request with title: %s, dueDate: %v", title, dueDateValue)

    task := &model.Task{
        Title:       title,
        Description: description,
        Status:      "pending",
        DueDate:     dueDate,
    }

    // Insert task into the database
    res, err := db.DB.Exec("INSERT INTO tasks (title, description, status, due_date) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, dueDateValue)
    if err != nil {
        log.Printf("Error inserting task: %v", err)
        return nil, err
    }

    // Get the last inserted ID
    id, err := res.LastInsertId()
    if err != nil {
        log.Printf("Error getting last insert ID: %v", err)
        return nil, err
    }
    task.ID = fmt.Sprintf("%d", id)

    log.Printf("Successfully added task with ID: %s, dueDate: %v", task.ID, dueDateValue)

    return task, nil
}

// UpdateTask is the resolver for the updateTask field.
func (r *mutationResolver) UpdateTask(ctx context.Context, id string, title *string, description *string, status *string, dueDate *string) (*model.Task, error) {
	var customTask model.CustomTask
	err := db.DB.Get(&customTask, "SELECT id, title, description, status, due_date FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	// Update fields if they are provided
	if title != nil {
		customTask.Title = *title
	}
	if description != nil {
		customTask.Description = description
	}
	if status != nil {
		customTask.Status = *status
	}
	if dueDate != nil {
		customTask.DueDate = dueDate
	}

	// Update task in the database
	_, err = db.DB.Exec("UPDATE tasks SET title=?, description=?, status=?, due_date=? WHERE id=?", customTask.Title, customTask.Description, customTask.Status, customTask.DueDate, id)
	if err != nil {
		return nil, err
	}

	task := &model.Task{
		ID:          customTask.ID,
		Title:       customTask.Title,
		Description: customTask.Description,
		Status:      customTask.Status,
		DueDate:     customTask.DueDate,
	}

	log.Printf("Successfully updated task with ID: %s, dueDate: %v", task.ID, task.DueDate)

	return task, nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	_, err := db.DB.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting task with ID: %s, error: %v", id, err)
		return false, err
	}

	log.Printf("Successfully deleted task with ID: %s", id)
	return true, nil
}

// MarkTaskCompleted is the resolver for the markTaskCompleted field.
func (r *mutationResolver) MarkTaskCompleted(ctx context.Context, id string) (*model.Task, error) {
	var customTask model.CustomTask
	err := db.DB.Get(&customTask, "SELECT id, title, description, status, due_date FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	// Update status to completed
	customTask.Status = "completed"
	_, err = db.DB.Exec("UPDATE tasks SET status=? WHERE id=?", customTask.Status, id)
	if err != nil {
		return nil, err
	}

	task := &model.Task{
		ID:          customTask.ID,
		Title:       customTask.Title,
		Description: customTask.Description,
		Status:      customTask.Status,
		DueDate:     customTask.DueDate,
	}

	log.Printf("Successfully marked task as completed with ID: %s", task.ID)

	return task, nil
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	var customTasks []*model.CustomTask
	err := db.DB.Select(&customTasks, "SELECT id, title, description, status, due_date FROM tasks")
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for _, ct := range customTasks {
		tasks = append(tasks, &model.Task{
			ID:          ct.ID,
			Title:       ct.Title,
			Description: ct.Description,
			Status:      ct.Status,
			DueDate:     ct.DueDate,
		})
	}

	// Debugging output to check dueDate values
	for _, task := range tasks {
		log.Printf("Task ID: %s, DueDate: %v", task.ID, task.DueDate)
	}

	return tasks, nil
}

// Task is the resolver for the task field.
func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {
	var customTask model.CustomTask
	err := db.DB.Get(&customTask, "SELECT id, title, description, status, due_date FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	// Convert customTask to model.Task
	task := &model.Task{
		ID:          customTask.ID,
		Title:       customTask.Title,
		Description: customTask.Description,
		Status:      customTask.Status,
		DueDate:     customTask.DueDate,
	}

	log.Printf("Retrieved task with ID: %s, DueDate: %v", task.ID, task.DueDate)

	return task, nil
}