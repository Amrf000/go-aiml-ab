package ab

import (
	"fmt"
)

type History struct {
	MHistory [maxHistory]interface{}
	MName    string
}

const maxHistory = 32
const unknownHistoryItem = "unknown"

func NewHistory() *History {
	return &History{
		MHistory: [maxHistory]interface{}{},
		MName:    "unknown",
	}
}

func NewHistoryWithName(name string) *History {
	return &History{
		MHistory: [maxHistory]interface{}{},
		MName:    name,
	}
}

func (h *History) Add(item interface{}) {
	for i := maxHistory - 1; i > 0; i-- {
		h.MHistory[i] = h.MHistory[i-1]
	}
	h.MHistory[0] = item
}

func (h *History) Get(index int) interface{} {
	if index < maxHistory {
		if h.MHistory[index] == nil {
			return nil
		}
		return h.MHistory[index]
	}
	return nil
}

func (h *History) GetString(index int) string {
	if index < maxHistory {
		if h.MHistory[index] == "" {
			return unknownHistoryItem
		}
		if str, ok := h.MHistory[index].(string); ok {
			return str
		}
		return fmt.Sprintf("%v", h.MHistory[index])
	}
	return ""
}

func (h *History) PrintHistory() {
	for i := 0; i < maxHistory; i++ {
		if item := h.Get(i); item != "" {
			if val, ok := item.(string); ok {
				fmt.Printf("%s History %d = %v\n", h.MName, i+1, val)
			}
			if _, ok := item.(*History); ok {
				item.(*History).PrintHistory()
			}
		}
	}
}
