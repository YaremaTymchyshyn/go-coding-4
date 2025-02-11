package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type BankAccount struct {
	Number  string
	Balance float64
	Credit  float64
	mu      sync.Mutex
}

func (acc *BankAccount) SetBalance(amount float64) {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	acc.Balance = amount
}

func (acc *BankAccount) GetBalance() float64 {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	return acc.Balance
}

func (acc *BankAccount) SetCredit(credit float64) {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	acc.Credit = credit
}

func (acc *BankAccount) GetCredit() float64 {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	return acc.Credit
}

type BankServer struct {
	accounts map[string]*BankAccount
	mu       sync.Mutex
}

func NewBankServer() *BankServer {
	return &BankServer{
		accounts: make(map[string]*BankAccount),
	}
}

func (s *BankServer) GetAccount(number string) *BankAccount {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.accounts[number]
}

func (s *BankServer) AddAccount(account *BankAccount) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[account.Number] = account
}

func handleConnection(conn net.Conn, server *BankServer) {
	defer conn.Close()

	var request map[string]interface{}
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		if err := decoder.Decode(&request); err != nil {
			fmt.Println("Error decoding request:", err)
			return
		}

		action := request["action"].(string)
		accountNumber := request["number"].(string)
		account := server.GetAccount(accountNumber)

		if account == nil {
			encoder.Encode(map[string]string{"error": "Account not found"})
			continue
		}

		switch action {
		case "set_balance":
			amount := request["amount"].(float64)
			account.SetBalance(amount)
			encoder.Encode(map[string]string{"status": "Balance updated"})

		case "get_balance":
			balance := account.GetBalance()
			encoder.Encode(map[string]float64{"balance": balance})

		case "set_credit":
			credit := request["credit"].(float64)
			account.SetCredit(credit)
			encoder.Encode(map[string]string{"status": "Credit updated"})

		case "get_credit":
			credit := account.GetCredit()
			encoder.Encode(map[string]float64{"credit": credit})

		default:
			encoder.Encode(map[string]string{"error": "Unknown action"})
		}
	}
}

func main() {
	server := NewBankServer()
	server.AddAccount(&BankAccount{Number: "12345", Balance: 1000, Credit: 500})
	server.AddAccount(&BankAccount{Number: "67890", Balance: 2000, Credit: 1000})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Bank server is running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, server)
	}
}
