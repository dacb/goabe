package life

import (
	"context"
	"fmt"
)

func (life *matrix) printMatrix(ctx context.Context) error {
	//time.Sleep(1 * time.Second)
	fmt.Print("\033[H\033[2J")
	for xi := 0; xi < life.x; xi++ {
		for yi := 0; yi < life.y; yi++ {
			c := '.'
			if life.mat[xi][yi].alive {
				c = 'X'
			}
			fmt.Printf(" %c", c)
		}
		fmt.Printf("\n")
	}
	return nil
}

func (life *matrix) saveMatrix(ctx context.Context, filename string) error {
	return nil
}

func (life *matrix) loadMatrix(ctx context.Context, filename string) error {
	return nil
}
