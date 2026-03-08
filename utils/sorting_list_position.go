package utils

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/google/uuid"
)

func SortListByPosition(lists []models.List, order []uuid.UUID) []models.List{

	if len(order) == 0 {
		return lists
	}

	ordered := make([]models.List, 0, len(order))

	listMap := make(map[uuid.UUID]models.List)

	for _, l := range lists {
		listMap[l.PublicID] = l
	}

	for _,id := range order {
		if l,ok :=listMap[id];ok {
			ordered = append(ordered, l)
		}
	}

	return ordered
}