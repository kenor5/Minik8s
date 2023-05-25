package uuid

import (
	"fmt"
	"github.com/google/uuid"
)

func UUID() string {
	uid := uuid.NewString()
	uid5 := uid[:5]
	fmt.Print(uid5)
	return uid
}
