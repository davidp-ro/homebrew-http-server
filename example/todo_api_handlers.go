package example

import (
	"encoding/json"
	"strconv"

	"github.com/davidp-ro/homebrew-http-server/server"
	"github.com/davidp-ro/homebrew-http-server/utils"
)

// Demo (in-memory "db") for todos
var todos []Todo = []Todo{
	{ID: GetRandomId(), Title: "Buy milk", Completed: false},
	{ID: GetRandomId(), Title: "Meet with John", Completed: false},
	{ID: GetRandomId(), Title: "Go to the gym", Completed: true},
}

func GetTodos(data server.HandlerData) []byte {
	todosJson, err := json.Marshal(todos)
	if err != nil {
		return data.Server.RespondWith(500, "Error marshalling todo list")
	}
	return data.Server.RespondWith(200, string(todosJson), map[string]string{
		"content-type": "application/json",
	})
}

func GetTodo(data server.HandlerData) []byte {
	id := data.PathParams["$id"]
	if id == "" {
		return data.Server.RespondWith(400, "Missing Todo ID")
	}

	todoId, err := strconv.Atoi(id)
	if err != nil {
		return data.Server.RespondWith(400, "Invalid Todo ID")
	}

	result := utils.Filter(todos, func(t Todo) bool {
		return t.ID == todoId
	})

	if len(result) == 0 {
		return data.Server.RespondWith(404, "Todo not found")
	} else if len(result) > 1 {
		return data.Server.RespondWith(500, "Multiple Todos with the same ID")
	}

	todoJson, err := TodoToJSON(result[0])
	if err != nil {
		return data.Server.RespondWith(500, "Error marshalling todo")
	}

	return data.Server.RespondWith(200, todoJson, map[string]string{
		"content-type": "application/json",
	})
}

func CreateTodo(data server.HandlerData) []byte {
	req := data.Request.Body.JSON

	todoTitle := req["title"]
	if todoTitle == "" {
		return data.Server.RespondWith(400, "Missing Todo title in body")
	}

	title, ok := todoTitle.(string)
	if !ok {
		return data.Server.RespondWith(400, "Invalid Todo title in body")
	}

	todos = append(todos, Todo{ID: GetRandomId(), Title: title, Completed: false})

	// Return the updated list of todos
	todosJson, err := json.Marshal(todos)
	if err != nil {
		return data.Server.RespondWith(500, "Error marshalling todo list")
	}

	return data.Server.RespondWith(201, string(todosJson), map[string]string{
		"content-type": "application/json",
	})
}

func DeleteTodo(data server.HandlerData) []byte {
	id := data.PathParams["$id"]
	if id == "" {
		return data.Server.RespondWith(400, "Missing Todo ID")
	}

	todoId, err := strconv.Atoi(id)
	if err != nil {
		return data.Server.RespondWith(400, "Invalid Todo ID")
	}

	// Remove the todo with the given ID
	todos = utils.Filter(todos, func(t Todo) bool {
		return t.ID != todoId
	})

	// Return the updated list of todos
	todosJson, err := json.Marshal(todos)
	if err != nil {
		return data.Server.RespondWith(500, "Error marshalling todo list")
	}

	return data.Server.RespondWith(200, string(todosJson), map[string]string{
		"content-type": "application/json",
	})
}
