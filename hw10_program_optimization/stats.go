package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)
	var user User

	decoder := jsoniter.NewDecoder(r)
	for decoder.More() {
		if err := decoder.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			emailPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			domainStat[emailPart]++
		}
	}

	return domainStat, nil
}
