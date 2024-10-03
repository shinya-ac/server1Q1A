package chatgpt

import (
	"strings"
)

func ParseQas(content string) ([]*Qas, error) {
	lines := strings.Split(content, "\n")
	var qasList []*Qas
	var currentQas *Qas

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Q.") {
			if currentQas != nil {
				qasList = append(qasList, currentQas)
			}
			currentQas = &Qas{Question: strings.TrimSpace(strings.TrimPrefix(line, "Q."))}
		} else if strings.HasPrefix(line, "A.") && currentQas != nil {
			currentQas.Answer = strings.TrimSpace(strings.TrimPrefix(line, "A."))
		}
	}
	if currentQas != nil {
		qasList = append(qasList, currentQas)
	}

	return qasList, nil
}
