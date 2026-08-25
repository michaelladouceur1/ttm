package models

import (
	"fmt"
	"strings"
)

// TaskSearch describes substring filters for task fields.
type TaskSearch struct {
	General     string
	Title       string
	Description string
	Tags        []string
	Priority    string
	Status      string
}

// ParseTaskSearch parses either a plain substring search or field filters such
// as "*tags:work,urgent *title:Task 1".
func ParseTaskSearch(input string) (TaskSearch, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return TaskSearch{}, fmt.Errorf("search query cannot be empty")
	}
	if !strings.Contains(input, "*") {
		return TaskSearch{General: input}, nil
	}

	var search TaskSearch
	for input != "" {

		input = strings.TrimSpace(input)
		if !strings.HasPrefix(input, "*") {
			return TaskSearch{}, fmt.Errorf("expected a field filter beginning with *")
		}

		separator := strings.IndexByte(input, ':')
		if separator < 2 {
			return TaskSearch{}, fmt.Errorf("invalid field filter %q", input)
		}
		field := input[1:separator]
		valueEnd := len(input)
		if next := strings.Index(input[separator+1:], " *"); next >= 0 {
			valueEnd = separator + 1 + next
		}
		value := strings.TrimSpace(input[separator+1 : valueEnd])
		if value == "" {
			return TaskSearch{}, fmt.Errorf("search value for %q cannot be empty", field)
		}

		switch field {
		case "title":
			search.Title = value
		case "description":
			search.Description = value
		case "tags":
			search.Tags = strings.Split(value, ",")
			for i := range search.Tags {
				search.Tags[i] = strings.TrimSpace(search.Tags[i])
				if search.Tags[i] == "" {
					return TaskSearch{}, fmt.Errorf("tag filters cannot be empty")
				}
			}
		case "priority":
			search.Priority = value
		case "status":
			search.Status = value
		default:
			return TaskSearch{}, fmt.Errorf("unsupported search field %q", field)
		}

		input = input[valueEnd:]
	}

	return search, nil
}
