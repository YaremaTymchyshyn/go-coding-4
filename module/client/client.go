package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	for {
		fmt.Print("Enter command (set_balance, get_balance, set_credit, get_credit) and account number: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		args := strings.Split(text, " ")

		if len(args) < 2 {
			fmt.Println("Invalid command")
			continue
		}

		action, accountNumber := args[0], args[1]
		request := map[string]interface{}{"action": action, "number": accountNumber}

		switch action {
		case "set_balance":
			if len(args) < 3 {
				fmt.Println("Missing amount")
				continue
			}
			amount, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			request["amount"] = amount

		case "set_credit":
			if len(args) < 3 {
				fmt.Println("Missing credit amount")
				continue
			}
			credit, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Println("Invalid credit amount")
				continue
			}
			request["credit"] = credit

		case "get_balance", "get_credit":

		default:
			fmt.Println("Unknown command")
			continue
		}

		if err := encoder.Encode(request); err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}

		var response map[string]interface{}
		if err := decoder.Decode(&response); err != nil {
			fmt.Println("Error receiving response:", err)
			continue
		}

		fmt.Println("Response:", response)
	}
}
