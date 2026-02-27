package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID


func (uuidArray *UUIDArray) Scan(value interface{}) error{
	var str string

	switch v := value.(type){
	case []byte: str = string(v)
	case string: str = v
	default:
		return errors.New("Failed to parse UUIArray: unsupported data type")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	parts := strings.Split(str, ",")

	*uuidArray = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`)) // akan menghapus spasi dan "
		if s == "" {
			continue
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID in Array : %v", err)
		}
		*uuidArray = append(*uuidArray, u)

	}

	return nil
}

func (uuidArray UUIDArray) Value()(driver.Value, error){
	if len(uuidArray) == 0 {
		return "{}", nil
	}

	postgresFormart := make([]string, 0, len(uuidArray))

	for _, value := range uuidArray{
		postgresFormart = append(postgresFormart, fmt.Sprintf(`"%v"`, value))
	}

	return "{"+ strings.Join(postgresFormart, ",") +"}", nil
}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}