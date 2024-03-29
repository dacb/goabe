package life

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/dacb/goabe/plugins"
	"github.com/spf13/viper"
)

type cell struct {
	alive     bool
	aliveNext bool
	neighbors []*cell
}

type matrix struct {
	x, y      int       // dimensions
	cells     []cell    // the matrix of cells allocated linearly
	mat       [][]*cell // a matrix to be addressed by the cell dimension, points to above
	chunkSize int       // the size of chunks in the matrix for each thread to process
}

var life matrix
var threads int

var rng *rand.Rand

func Register() {
	plugins.LoadedPlugins = append(plugins.LoadedPlugins, plugins.Plugin{Init, Name, Version, Description, GetHooks, PreRun, PostRun})
}

// main initiailization function for the plugin
func Init(ctx context.Context) error {
	log, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		return errors.New("no logger found on the current context")
	}

	threadCount, ok := ctx.Value("threads").(int)
	if !ok {
		return errors.New("missing number of threads in current context")
	}
	threads = threadCount

	// initialize the random seed
	random_seed := viper.GetInt64("random_seed")
	rng = rand.New(rand.NewSource((random_seed)))

	log.Info(fmt.Sprintf("Life plugin Init function was called for %d threads w/ %d as random_seed", threads, random_seed))

	// initialize the data structures for the module
	//viper.SetDefault("life.in_filename", "IN.LIF")
	viper.SetDefault("life.out_filename", "OUT.LIF")

	viper.SetDefault("life.x_size", 16)
	viper.SetDefault("life.y_size", 32)
	life.x = viper.GetInt("life.x_size")
	life.y = viper.GetInt("life.y_size")

	if life.x < 3 || life.y < 3 {
		log.Error("minimum size of matrix must be 3 x 3")
		return errors.New("minimum size of matrix must be 3 x 3")
	}

	// allocate the cellular matrix
	life.cells = make([]cell, life.x*life.y)
	idx := 0
	life.mat = make([][]*cell, life.x)
	for xi := 0; xi < life.x; xi++ {
		life.mat[xi] = make([]*cell, life.y)
		for yi := 0; yi < life.y; yi++ {
			life.mat[xi][yi] = &life.cells[idx]
			idx += 1
		}
	}
	// now setup the neighbors
	idx = 0
	for xi := 0; xi < life.x; xi++ {
		for yi := 0; yi < life.y; yi++ {
			// loop over the matrix neighbors and setup the neighbor list
			// do this using periodic boundaries; doing this once here
			// allows us to avoid all the ifs when iterating over the neighbors
			// at the cost of a pointer dereference, should be a good deal
			// performance wise
			life.cells[idx].neighbors = make([]*cell, 8)
			nidx := 0
			for nxi := -1; nxi <= 1; nxi++ {
				for nyi := -1; nyi <= 1; nyi++ {
					// skip outselves as a neighbor
					if nxi == 0 && nyi == 0 {
						continue
					}
					// calculate neihbor index w/ periodic boundaries
					nxiPB := nxi + xi
					if nxiPB < 0 {
						nxiPB = life.x - 1
					} else if nxiPB >= life.x {
						nxiPB = 0
					}
					nyiPB := nyi + yi
					if nyiPB < 0 {
						nyiPB = life.y - 1
					} else if nyiPB >= life.y {
						nyiPB = 0
					}
					// set the neighbor in update the index
					life.cells[idx].neighbors[nidx] = life.mat[nxiPB][nyiPB]
					nidx += 1
				}
			}
			idx += 1
		}
	}

	// initialize the matrix states
	filename := viper.GetString("life.in_filename")
	if filename != "" {
		err := life.loadMatrix(ctx, filename)
		if err != nil {
			log.Error(fmt.Sprintf("unable to load starting matrix file file '%s'", filename))
			return err
		}
		life.printMatrix(ctx)
	} else {
		aliveCells := 0
		for idx := 0; idx < life.x*life.y; idx++ {
			life.cells[idx].alive = rng.Uint32()&(1<<31) == 0 // random true false from integer
			if life.cells[idx].alive {
				aliveCells += 1
			}
		}
		log.Info(fmt.Sprintf("there are %d alive cells at the start", aliveCells))
	}

	return nil
}

// major, minor, patch
func Version() (int, int, int) {
	return 0, 1, 0
}

// returns the short name of the module as a string
func Name() string {
	return "Life"
}

// returns a short description of the module as a string
func Description() string {
	return "Conway's game of plugin for code template"
}

func GetHooks() []plugins.Hook {
	var hooks []plugins.Hook
	hooks = append(hooks, plugins.Hook{0, nil, ThreadSubStep0, "thread calculate next state"})
	hooks = append(hooks, plugins.Hook{1, CoreSubStep1, nil, "core update next state"})

	return hooks
}

// do before each set of steps
func PreRun(ctx context.Context) error {
	//log := ctx.Value("log").(*slog.Logger).With("plugin", Name())

	return nil
}

// do after each set of steps
func PostRun(ctx context.Context) error {
	//log := ctx.Value("log").(*slog.Logger).With("plugin", Name())

	filename := viper.GetString("life.out_filename")
	life.saveMatrix(ctx, filename)

	return nil
}

// note this logs through the context
func CoreSubStep1(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger).With("plugin", Name())
	aliveCells := 0
	for idx := 0; idx < life.x*life.y; idx++ {
		life.cells[idx].alive = life.cells[idx].aliveNext
		if life.cells[idx].alive {
			aliveCells += 1
		}
	}
	log.Info(fmt.Sprintf("%d alive cells", aliveCells))
	life.printMatrix(ctx)
	return nil
}

// note this logs through the context
func ThreadSubStep0(ctx context.Context, id int, name string) error {
	//log := ctx.Value("log").(*slog.Logger).With("plugin", Name())

	// determine the chunksize
	chunkSize := (life.x * life.y) / threads
	if (life.x*life.y)%threads != 0 {
		chunkSize += 1
	}
	// iterative over this thread's chunk
	for idx := chunkSize * id; idx <= chunkSize*(id+1) && idx < life.x*life.y; idx++ {
		alive := 0
		for nidx := 0; nidx < 8; nidx++ {
			if life.cells[idx].neighbors[nidx].alive {
				alive += 1
			}
		}
		life.cells[idx].aliveNext = false
		if life.cells[idx].alive {
			if alive == 2 || alive == 3 {
				life.cells[idx].aliveNext = true
			}
		} else {
			if alive == 3 {
				life.cells[idx].aliveNext = true
			}
		}
	}

	return nil
}
