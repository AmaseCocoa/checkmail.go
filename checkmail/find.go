package main;

func CheckDomainExists(domain string) (bool, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	var exists int
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM domains WHERE domain = ?)", domain).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
