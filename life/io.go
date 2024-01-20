package life

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

func (life *matrix) printMatrix(ctx context.Context) error {
	//time.Sleep(1 * time.Second)
	fmt.Print("\033[H\033[2J")
	for yi := 0; yi < life.y; yi++ {
		for xi := 0; xi < life.x; xi++ {
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

	major, minor, build := Version()
	file.WriteString(fmt.Sprintf("#C goabe life plugin v%d.%d.%d\n", major, minor, build))
	file.WriteString(fmt.Sprintf("x = %d, y = %d, rule = B3/S23\n", life.x, life.y))
	var c = '-' // unknown state, b = dead, o = alive
	for y := 0; y < life.y; y++ {
		n := 0
		c = '-'
		for x := 0; x < life.x; x++ {
			if life.mat[x][y].alive {
				if c == '-' {
					n = 1
					c = 'o'
				} else if c == 'b' {
					if n > 1 {
						file.WriteString(fmt.Sprintf("%d%c", n, c))
					} else {
						file.WriteString(fmt.Sprintf("%c", c))
					}
					c = 'o'
					n = 1
				} else {
					n = n + 1
				}
			} else {
				if c == '-' {
					n = 1
					c = 'b'
				} else if c == 'o' {
					if n > 1 {
						file.WriteString(fmt.Sprintf("%d%c", n, c))
					} else {
						file.WriteString(fmt.Sprintf("%c", c))
					}
					c = 'b'
					n = 1
				} else {
					n = n + 1
				}
			}
		}
		if c == 'o' {
			if n > 1 {
				file.WriteString(fmt.Sprintf("%d%c", n, c))
			} else {
				file.WriteString(fmt.Sprintf("%c", c))
			}
		}
		if y != life.y-1 {
			file.WriteString("$\n")
		} else {
			file.WriteString("!\n")
		}
	}

	return nil
}

type parserState int

const (
	start parserState = iota
	formatKnown
	done
)

// makes no assumptions about the matrix being empty
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
	var x, y, B, S int // x, y dimensions of pattern, B and S are rule values
	xi := 0
	yi := 0
	for scanner.Scan() {
		log := log.With("matrix_file", filename)
		text := scanner.Text()
		words := strings.Fields(text)
		switch state {
		case done:
			break
		case start:
			if len(words) < 1 {
				log.Warn(fmt.Sprintf("line %d: expecting to find a non-empty line", line))
			} else if words[0] == "#C" {
				log.Info(text)
			} else {
				parsed, err := fmt.Sscanf(text, "x = %d, y = %d, rule = B%d/S%d", &x, &y, &B, &S)
				if parsed == 2 || parsed == 4 {
					state = formatKnown
				}
				if parsed == 4 {
					// we only support patterns from standard rules, error out
					if B != 3 && S != 23 {
						error_msg := fmt.Sprintf("line %d: unsupported ruleset for life", line)
						log.Error(error_msg)
						return errors.New("unsupported data in RLE file")
					}
				}
				if x < 0 || x > life.x || y < 0 || y > life.y {
					error_msg := fmt.Sprintf("line %d: pattern dimensions exceed matrix size (%d by %d)", line, x, y)
					log.Error(error_msg)
					return errors.New("pattern too large for current run")
				}
				if state == start && err != nil {
					log.Error(fmt.Sprintf("line %d: an error occurred parsing the pattern format line", line))
				} else {
					log.Info(fmt.Sprintf("input pattern is %d by %d and loaded without a problem", x, y))
				}
			}
		case formatKnown:
			if len(words) < 1 {
				log.Error(fmt.Sprintf("line %d: line too short", line))
			}
			// to parse these lines, lets split the entire string
			// then we will join tokens that are numbers until we
			// have only alternating numbers and letters or just
			// letters
			re := regexp.MustCompile("([0-9])*(.){1}")
			split := re.FindAllString(text, -1)
			//fmt.Printf("found %d words in text: %s\n", len(split), text)
			//fmt.Printf("%q\n", split)
			for i := 0; i < len(split); i++ {
				//fmt.Printf("word: %s\n", split[i])
				var n int
				parsed, _ := fmt.Sscanf(split[i], "%d", &n)
				if parsed != 1 {
					n = 1
				}
				// last character of this word is the type
				c := split[i][len(split[i])-1]
				for ni := 0; ni < n; ni++ {
					switch c {
					case 'o':
						if xi < 0 || xi >= life.x || yi < 0 || yi >= life.y {
							msg := fmt.Sprintf("line %d: x (%d) or y (%d) dimensions out of range (%d, %d) when parsing RLE data\n", line, xi, yi, life.x, life.y)
							log.Error(msg)
							return errors.New("unable to parse the pattern from the file due to invalid sizes")
						}
						life.mat[xi][yi].alive = true
						xi = xi + 1
					case 'b':
						if xi < 0 || xi >= life.x || yi < 0 || yi >= life.y {
							msg := fmt.Sprintf("line %d: x (%d) or y (%d) dimensions out of range (%d, %d) when parsing RLE data\n", line, xi, yi, life.x, life.y)
							log.Error(msg)
							return errors.New("unable to parse the pattern from the file due to invalid sizes")
						}
						life.mat[xi][yi].alive = false
						xi = xi + 1
					case '$':
						yi = yi + 1
						xi = 0
					case '!':
						state = done
						break
					default:
						msg := fmt.Sprintf("line %d: unexpected character found '%c'\n", line, c)
						log.Error(msg)
						return errors.New("unable to parse the pattern from the file")
					}
				}
			}
		}
		line += 1
	}
	if state == start {
		msg := fmt.Sprintf("no valid matrix found in file '%s'", filename)
		log.Error(msg)
		return errors.New("unable to parse input file")
	}

	if err := scanner.Err(); err != nil {
		log.Error(fmt.Sprintf("an error occurred reading from the file '%s'", filename))
		return err
	}

	return nil
}
