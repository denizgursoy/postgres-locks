package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/denizgursoy/postgres-locks/db"
)

func main() {
	// prepare connections
	connection := db.GetDB(context.Background())
	defer connection.Close()
	var orders []db.Order

	// first transaction
	firstTransaction, _ := connection.Beginx()
	fmt.Println("Locking the orders of customer 1 for another select query")
	if err := firstTransaction.Select(&orders, `SELECT * FROM orders WHERE customer_id=1 FOR UPDATE`); err != nil {
		log.Fatalf("Error locking the orders the customer 1: %v", err)
	}
	fmt.Printf("%d record are selected in the first select for update query\n", len(orders))

	fmt.Println("First select for update query will be commited 5 seconds later")
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("First select for update query is committed")
		if err := firstTransaction.Commit(); err != nil {
			log.Fatalf("Error committing the select for update: %v", err)
		}
	}()

	fmt.Println("Another transaction is trying get the locks on the same records")
	fmt.Println("It will try to select and records are locked therefore it will not wait if the records are locked")
	if err := connection.Select(&orders, `SELECT * FROM orders WHERE customer_id=1 FOR UPDATE NOWAIT`); err != nil {
		log.Fatalf("Error locking the orders the customer 1 from another transaction: %v", err)
	}
}
