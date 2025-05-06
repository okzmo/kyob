package utils

import (
	"context"
	"math/rand"

	"github.com/okzmo/kyob/db"
)

type Coordinates struct {
	X int
	Y int
}

func GenerateRandomCoordinates() (int32, int32, error) {
	var x, y int32
	for {
		x = rand.Int31n(2000)
		y = rand.Int31n(2000)
		res, err := db.Query.CheckServerPosition(context.TODO(), db.CheckServerPositionParams{
			Column1: x,
			Column2: y,
		})
		if err != nil {
			return 0, 0, err
		}

		if res.RowsAffected() == 0 {
			goto Found
		}
	}

Found:

	return x, y, nil
}
