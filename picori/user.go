package picori

import "time"

type User struct {
    Id string
    Username string
    Created time.Time
}

type UserFilter struct {
    Id string
    Username string
}
