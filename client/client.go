package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"grpc-redis/protos/todo/protos/todo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := todo.NewTodoServiceClient(conn)

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <command> [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) != 4 {
			fmt.Println("Usage: client add <title> <description>")
			return
		}
		title := os.Args[2]
		description := os.Args[3]
		addTodo(client, title, description)
	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: client get <id>")
			return
		}
		id := os.Args[2]
		getTodo(client, id)
	case "update":
		if len(os.Args) != 5 {
			fmt.Println("Usage: client update <id> <title> <description>")
			return
		}
		id := os.Args[2]
		title := os.Args[3]
		description := os.Args[4]
		updateTodo(client, id, title, description)
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: client delete <id>")
			return
		}
		id := os.Args[2]
		deleteTodo(client, id)
	case "list":
		if len(os.Args) != 2 {
			fmt.Println("Usage: client list")
			return
		}
		listTodos(client)
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage: client <command> [arguments]")
	}
}

func addTodo(client todo.TodoServiceClient, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.AddTodo(ctx, &todo.AddTodoRequest{Title: title, Description: description})
	if err != nil {
		log.Fatalf("could not add todo: %v", err)
	}
	fmt.Printf("Todo added with ID: %s\n", resp.Id)
}

func getTodo(client todo.TodoServiceClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetTodo(ctx, &todo.GetTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not get todo: %v", err)
	}
	fmt.Printf("Todo: ID: %s, Title: %s, Description: %s\n", resp.Id, resp.Title, resp.Description)
}

func updateTodo(client todo.TodoServiceClient, id, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.UpdateTodo(ctx, &todo.UpdateTodoRequest{Id: id, Title: title, Description: description})
	if err != nil {
		log.Fatalf("could not update todo: %v", err)
	}
	fmt.Printf("Todo updated: %s\n", resp.Success)
}

func deleteTodo(client todo.TodoServiceClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.DeleteTodo(ctx, &todo.DeleteTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not delete todo: %v", err)
	}
	fmt.Printf("Todo deleted: %s\n", resp.Success)
}

func listTodos(client todo.TodoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.ListTodos(ctx, &todo.ListTodosRequest{})
	if err != nil {
		log.Fatalf("could not list todos: %v", err)
	}
	fmt.Println("Todos:")
	for _, todoItem := range resp.Todos {
		fmt.Printf("ID: %s, Title: %s, Description: %s\n", todoItem.Id, todoItem.Title, todoItem.Description)
	}
}
