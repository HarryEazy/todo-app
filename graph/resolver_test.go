package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_app/db"
	"todo_app/graph/generated"
	"todo_app/graph/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
)

func setup() *handler.Server {
    // Initialize the database
    db.InitDB()

    // Clear the database before each test to avoid state conflicts
    db.DB.MustExec("DELETE FROM tasks")

    // Initialize the GraphQL server
    resolver := &Resolver{}
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

    return srv
}


// Helper function to send GraphQL requests
func sendGraphQLRequest(srv *handler.Server, query string) (*http.Response, error) {
	reqBody := map[string]string{
		"query": query,
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/query", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	return w.Result(), nil
}

// TestAddTask tests the AddTask mutation
func TestAddTask(t *testing.T) {
	srv := setup()

	query := `
		mutation {
			addTask(title: "Test Task", description: "This is a test task") {
				id
				title
				description
				status
			}
		}
	`

	resp, err := sendGraphQLRequest(srv, query)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Data struct {
			AddTask model.Task
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Assert the task is added correctly
	assert.NotEmpty(t, result.Data.AddTask.ID)
	assert.Equal(t, "Test Task", result.Data.AddTask.Title)
	assert.Equal(t, "This is a test task", *result.Data.AddTask.Description)
	assert.Equal(t, "pending", result.Data.AddTask.Status)
}

// TestTasks tests the Tasks query
func TestTasks(t *testing.T) {
	srv := setup()

	// First, add a task to test the retrieval
	_ = db.DB.MustExec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", "Task 1", "Description 1", "pending")

	query := `
		query {
			tasks {
				id
				title
				description
				status
			}
		}
	`

	resp, err := sendGraphQLRequest(srv, query)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Tasks []model.Task
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Assert that tasks are retrieved correctly
	assert.NotEmpty(t, result.Data.Tasks)
	assert.Equal(t, "Task 1", result.Data.Tasks[0].Title)
	assert.Equal(t, "Description 1", *result.Data.Tasks[0].Description)
	assert.Equal(t, "pending", result.Data.Tasks[0].Status)
}

// TestUpdateTask tests the UpdateTask mutation
func TestUpdateTask(t *testing.T) {
	srv := setup()

	// First, add a task to update
	result := db.DB.MustExec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", "Task 1", "Description 1", "pending")
	id, _ := result.LastInsertId()

	query := fmt.Sprintf(`
		mutation {
			updateTask(id: "%d", title: "Updated Task", status: "completed") {
				id
				title
				description
				status
			}
		}
	`, id)

	resp, err := sendGraphQLRequest(srv, query)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var updateResult struct {
		Data struct {
			UpdateTask model.Task
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&updateResult)
	assert.NoError(t, err)

	// Assert that the task is updated correctly
	assert.Equal(t, fmt.Sprintf("%d", id), updateResult.Data.UpdateTask.ID)
	assert.Equal(t, "Updated Task", updateResult.Data.UpdateTask.Title)
	assert.Equal(t, "completed", updateResult.Data.UpdateTask.Status)
}

// TestDeleteTask tests the DeleteTask mutation
func TestDeleteTask(t *testing.T) {
	srv := setup()

	// First, add a task to delete
	result := db.DB.MustExec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", "Task to Delete", "Description", "pending")
	id, _ := result.LastInsertId()

	query := fmt.Sprintf(`
		mutation {
			deleteTask(id: "%d")
		}
	`, id)

	resp, err := sendGraphQLRequest(srv, query)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var deleteResult struct {
		Data struct {
			DeleteTask bool
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&deleteResult)
	assert.NoError(t, err)

	// Assert that the task is deleted
	assert.True(t, deleteResult.Data.DeleteTask)

	// Verify task deletion from the database
	var count int
	db.DB.Get(&count, "SELECT COUNT(*) FROM tasks WHERE id=?", id)
	assert.Equal(t, 0, count)
}

// TestMarkTaskCompleted tests the MarkTaskCompleted mutation
func TestMarkTaskCompleted(t *testing.T) {
	srv := setup()

	// First, add a task to mark as completed
	result := db.DB.MustExec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", "Task to Complete", "Description", "pending")
	id, _ := result.LastInsertId()

	query := fmt.Sprintf(`
		mutation {
			markTaskCompleted(id: "%d") {
				id
				title
				description
				status
			}
		}
	`, id)

	resp, err := sendGraphQLRequest(srv, query)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var completeResult struct {
		Data struct {
			MarkTaskCompleted model.Task
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&completeResult)
	assert.NoError(t, err)

	// Assert that the task is marked as completed
	assert.Equal(t, fmt.Sprintf("%d", id), completeResult.Data.MarkTaskCompleted.ID)
	assert.Equal(t, "completed", completeResult.Data.MarkTaskCompleted.Status)
}