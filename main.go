package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Account struct {
	Name    string
	Balance int
	Owe     map[string]Owe
}

type Owe struct {
	Name   string
	Amount int
	Type   int // 1: owed to, 2: owed from
}

var (
	userAccount = map[string]Account{}
)

var loginUser Account

func main() {
	input := command()
	action(input)
}

func command() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(scanner.Text(), " ")
}

func action(input []string) {
	switch input[0] {
	case "login":
		if _, ok := userAccount[input[1]]; !ok {
			userAccount[input[1]] = Account{
				Name: input[1],
			}
		}
		loginUser = userAccount[input[1]]
		fmt.Printf("Hello %s", loginUser.Name)
		fmt.Println()
		fmt.Printf("Your balance is $%d", loginUser.Balance)
		fmt.Println()
		if len(loginUser.Owe) > 0 {
			for user, data := range loginUser.Owe {
				if data.Type == 1 {
					fmt.Printf("Owed $%d to %s", data.Amount, user)
					fmt.Println()
				}
				if data.Type == 2 {
					fmt.Printf("Owed $%d from %s", data.Amount, user)
					fmt.Println()
				}
			}
		}
		main()
	case "deposit":
		amount, err := strconv.Atoi(input[1])
		if err != nil {
			return
		}
		loginUser.Balance += amount
		userAccount[loginUser.Name] = loginUser
		fmt.Printf("Your balance is $%d", userAccount[loginUser.Name].Balance)
		fmt.Println()
		main()
	case "withdraw":
		amount, err := strconv.Atoi(input[1])
		if err != nil {
			return
		}
		loginUser.Balance -= amount
		userAccount[loginUser.Name] = loginUser
		fmt.Printf("Your balance is $%d", userAccount[loginUser.Name].Balance)
		fmt.Println()
		main()
	case "transfer":
		amount, err := strconv.Atoi(input[2])
		if err != nil {
			return
		}
		targetUser := Account{
			Name:    input[1],
			Balance: userAccount[input[1]].Balance,
		}
		if amount > loginUser.Balance {
			targetUser.Balance += amount - loginUser.Balance
			if len(loginUser.Owe) == 0 {
				loginUser.Owe = map[string]Owe{}
			}
			loginUser.Owe[targetUser.Name] = Owe{
				Name:   targetUser.Name,
				Amount: amount - loginUser.Balance,
				Type:   1,
			}
			if len(targetUser.Owe) == 0 {
				targetUser.Owe = map[string]Owe{}
			}
			targetUser.Owe[loginUser.Name] = Owe{
				Name:   loginUser.Name,
				Amount: amount - loginUser.Balance,
				Type:   2,
			}
			fmt.Printf("Transferred $%d to %s", loginUser.Balance, targetUser.Name)
			fmt.Println()
			fmt.Printf("Your balance is $%d", loginUser.Balance)
			fmt.Println()
			fmt.Printf("Owed $%d to %s", amount-loginUser.Balance, targetUser.Name)
			fmt.Println()
			loginUser.Balance = 0
		} else {
			loginUser.Balance -= amount
			targetUser.Balance += amount
			fmt.Printf("Transferred $%d to %s", amount, targetUser.Name)
			fmt.Println()
			fmt.Printf("Your balance is $%d", loginUser.Balance)
			fmt.Println()
		}
		userAccount[loginUser.Name] = loginUser
		userAccount[targetUser.Name] = targetUser
		main()
	case "logout":
		loginUser = Account{}
		main()
	}
}
