package life

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
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
	log := ctx.Value("log").(*slog.Logger)

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Error(fmt.Sprintf("unable to open output matrix to file '%s'", filename))
		return err
	}

	w := bufio.NewWriter(file)
	w.WriteString("#Life 1.05\n")
	major, minor, build := Version()
	w.WriteString(fmt.Sprintf("#D goabe life plugin v%d.%d.%d\n", major, minor, build))

	return nil
}

type parserState int

const (
	start parserState = iota
	formatKnown
	cellBlock
)

func (life *matrix) loadMatrix(ctx context.Context, filename string) error {
	log := ctx.Value("log").(*slog.Logger)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Error(fmt.Sprintf("unable to open input matrix from file '%s'", filename))
		return err
	}

	scanner := bufio.NewScanner(file)

	state := start
	line := 1
	var cellBlockX, cellBlockY int
	for scanner.Scan() {
		log := log.With("matrix_file", filename)
		text := scanner.Text()
		words := strings.Fields(text)
		switch state {
		case start:
			if text != "#Life 1.05" || len(words) < 1 {
				log.Error(fmt.Sprintf("line %d: expecting to find a file format line, i.e., #Life 1.05", line))
			}
			state = formatKnown
		case formatKnown:
			if len(words) < 1 {
				log.Error(fmt.Sprintf("file '%s', line %d: line too short", filename, line))
			}
			switch words[0] {
			case "#D":
				//fmt.Println("description line")
			case "#P":
				//fmt.Println("switch to cell block")
				if len(words) != 3 {
					log.Error(fmt.Sprintf("file '%s', line %d: unable to parse cell block: %s", filename, line, text))
				}
				parsed, err := fmt.Sscanf(text, "#P %d %d", &cellBlockX, &cellBlockY)
				if parsed != 2 || err != nil {
					log.Error(fmt.Sprint("file '%s', line %d: unable to parse dimensions of cell block: %s", filename, line, text))
				}

				state = cellBlock
			case "#N":
				//fmt.Println("normal rules")
			case "#R":
				if text != "#R 23/3" {
					log.Error(fmt.Sprintf("file '%s', line %d: only normal Conway rules are supported", filename, line))
				}
				//fmt.Println("rules section")

			default:
				log.Error(fmt.Sprintf("file '%s', line %d: unable to read line", filename, line))
			}
		case cellBlock:
			if len(words) < 1 {
				log.Error(fmt.Sprintf("file '%s', line %d: line too short", filename, line))
			}
			switch words[0] {
			case "#P":
				if len(words) != 3 {
					log.Error(fmt.Sprintf("file '%s', line %d: unable to parse cell block: %s", filename, line, text))
				}
				parsed, err := fmt.Sscanf(text, "#P %d %d", &cellBlockX, &cellBlockY)
				if parsed != 2 || err != nil {
					log.Error(fmt.Sprint("file '%s', line %d: unable to parse dimensions of cell block: %s", filename, line, text))
				}
			default:
				if len(words) != 1 {
					log.Error(fmt.Sprintf("file '%s', line %d: unexpected number of words on line: %s", filename, line, len(words)))
				}
				for i, c := range text {
					if c == '*' {
						x := cellBlockX + i + (life.x / 2)
						y := cellBlockY + (life.y / 2)
						if x < 0 || x >= life.x {
							log.Error(fmt.Sprintf("file '%s', line %d, pos %d: x (%d) is out of bounds [%d, %d)", filename, line, i+1, x, 0, life.x))
						} else if y < 0 || y > life.y {
							log.Error(fmt.Sprintf("file '%s', line %d, pos %d: y (%d) is out of bounds [%d, %d)", filename, line, i+1, y, 0, life.y))
						} else {
							life.mat[x][y].alive = true
						}
					}
				}
				cellBlockY = cellBlockY + 1
			}

		default:
			log.Error("unknown parser state")
		}
		line += 1
	}
	if state == start {
		msg := fmt.Sprintf("no valid matrix found in file '%s'", filename)
		log.Error(msg)
		return errors.New(msg)
	}

	if err := scanner.Err(); err != nil {
		log.Error(fmt.Sprintf("an error occurred reading from the file '%s'", filename))
		return err
	}

	return nil
}
