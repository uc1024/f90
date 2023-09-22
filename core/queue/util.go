package queue

import "strings"

func generateName(push []Pusher) string {
	names := []string{}
	for _, v := range push {
		names = append(names, v.Name())
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, ",")
}
