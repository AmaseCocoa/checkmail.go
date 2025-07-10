package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	dbMu sync.Mutex
)

func updateDomains(newDomains map[string]struct{}) error {
	existingDomains := make(map[string]struct{})

	rows, err := db.Query("SELECT domain FROM domains")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return err
		}
		existingDomains[domain] = struct{}{}
	}

	for domain := range newDomains {
		_, err := db.Exec("INSERT OR IGNORE INTO domains (domain) VALUES (?)", domain)
		if err != nil {
			return err
		}
	}

	for domain := range existingDomains {
		if _, exists := newDomains[domain]; !exists {
			_, err := db.Exec("DELETE FROM domains WHERE domain = ?", domain)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeEmptyKeys(iterable []string) []string {
	var result []string
	for _, s := range iterable {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

func checkEmailDomains() {
	for {
		list1, err := fetchList("https://raw.githubusercontent.com/disposable-email-domains/disposable-email-domains/refs/heads/main/disposable_email_blocklist.conf")
		if err != nil {
			fmt.Println("Error fetching list1:", err)
			continue
		}

		list2, err := fetchList("https://raw.githubusercontent.com/chan-mai/kukulu-disposable-email-list/refs/heads/main/domains.txt")
		if err != nil {
			fmt.Println("Error fetching list2:", err)
			continue
		}

		domains := make(map[string]struct{})
		for _, domain := range removeEmptyKeys(append(list1, list2...)) {
			domains[domain] = struct{}{}
		}

		err = updateDomains(domains)
		if err != nil {
			fmt.Println("Error updating domains:", err)
			continue
		}

		fmt.Println("Domain Updated :D")
		time.Sleep(1 * time.Hour)
	}
}

func fetchList(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}

func InitializeDB() {
	var err error
	db, err = sql.Open("sqlite", "file:domains.db?cache=shared")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS domains (domain TEXT PRIMARY KEY)")
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	go checkEmailDomains()
}
