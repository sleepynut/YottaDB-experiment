package custom

import "fmt"

type RecordNotFound struct {
	TbName string
	Id     string
}

func (r RecordNotFound) Error() string {
	return fmt.Sprintf("Record Not found - (%s): %s", r.TbName, r.Id)
}
