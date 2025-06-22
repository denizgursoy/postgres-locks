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
	fmt.Println("Locking the orders of customer 1 for UPDATE and DELETE")
	if err := firstTransaction.Select(&orders, `SELECT * FROM orders WHERE customer_id=1 FOR UPDATE`); err != nil {
		log.Fatalf("Error locking the orders the customer 1: %v", err)
	}
	fmt.Printf("%d orders are selected\n", len(orders))

	fmt.Println("Select for update query will be commited 5 seconds later")
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Select for update query is committed")
		if err := firstTransaction.Commit(); err != nil {
			log.Fatalf("Error committing the select for update: %v", err)
		}
	}()

	secondTransaction, _ := connection.Beginx()
	fmt.Println("Another transaction is trying to update the statuses of order of customer 1")
	fmt.Println("It will wait until the locks are removed")
	if err := secondTransaction.Select(&orders, `UPDATE orders SET status='DELIVERED' WHERE customer_id= 1`); err != nil {
		log.Fatalf("Error updating the locked row from another transaction: %v", err)
	}
	fmt.Println("Second query is executed and row are updated")
	if err := secondTransaction.Commit(); err != nil {
		log.Fatalf("Error committing another transaction: %v", err)
	}
}
