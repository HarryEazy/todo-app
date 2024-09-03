package model

// CustomTask is a custom struct for database interactions
type CustomTask struct {
    ID          string  `json:"id" db:"id"`
    Title       string  `json:"title" db:"title"`
    Description *string `json:"description,omitempty" db:"description"`
    Status      string  `json:"status" db:"status"`
    DueDate     *string `json:"dueDate,omitempty" db:"due_date"` // Use db tag to map to due_date
}
