package entity

import (
	"fmt"
	"time"
)

// AuthContext holds contextual information for authorization
type AuthContext struct {
	Time     time.Time
	IP       string
	Location string
	// Add other context attributes as needed
}

func (c AuthContext) Hash() string {
	return fmt.Sprintf("%d:%s:%s", c.Time.Unix(), c.IP, c.Location)
}

func (c *AuthContext) GetAttributes() map[string]interface{} {
	return map[string]interface{}{
		"time":     c.Time,
		"ip":       c.IP,
		"location": c.Location,
	}
}
