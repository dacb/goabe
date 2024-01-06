package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/dacb/goabe/plugins"
)

var pluginFilename string // when this is populated, the plugin has alread been initialized
var log *slog.Logger

// the empty struct used for the plugin
type plugin struct {
}

var PlugIn plugin

type cell struct {
	occupied      bool
	occupied_next bool
	neighbors     []*cell
}

type matrix struct {
	x, y      int       // dimensions
	cells     []cell    // the matrix of cells allocated linearly
	mat       [][]*cell // a matrix to be addressed by the cell dimension, points to above
	chunkSize int       // the size of chunks in the matrix for each thread to process
}

var life matrix

// main initiailization function for the plugin
func (p *plugin) Init(ctx context.Context, pluginFname string) error {
	mylog, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		return errors.New("no logger found on the current context")
	}
	log = mylog
	log.Info("Life plugin Init function was called")

	if pluginFilename != "" {
		log.Error("plugin has already been initialized? refusing to load the plugin twice")
		return errors.New("this plug in has already been loaded!")
	}
	pluginFilename = pluginFname

	// initialize the data structures for the module
	life.x = 42
	life.y = 32

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
		life.mat[xi] = make([]*cell, life.y)
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
					nxiPB := nxi
					if nxi < 0 {
						nxiPB = life.x - 1
					} else if nxi >= life.x {
						nxiPB = 0
					}
					nyiPB := nyi
					if nyi < 0 {
						nyiPB = life.y - 1
					} else if nyi >= life.y {
						nyiPB = 0
					}
					// set the neighbor in update the index
					life.cells[idx].neighbors[nidx] = life.mat[nxiPB][nyiPB]
					nidx += 1
				}
			}
		}
	}

	// initialize the matrix states
	occupiedCells := 0
	for _, cell := range life.cells {
		cell.occupied = rand.Uint32()&(1<<31) == 0 // random true false from integer
		if cell.occupied {
			occupiedCells += 1
		}
	}
	log.Info(fmt.Sprintf("there are %d alive cells at the start", occupiedCells))

	return nil
}

// major, minor, patch
func (p *plugin) Version() (int, int, int) {
	return 0, 1, 0
}

// returns the short name of the module as a string
func (p *plugin) Name() string {
	return "Life"
}

// returns a short description of the module as a string
func (p *plugin) Description() string {
	return "Conway's game of plugin for code template"
}

func (p *plugin) GetHooks() []plugins.Hook {
	var hooks []plugins.Hook
	hooks = append(hooks, plugins.Hook{0, nil, ThreadSubStep0, "thread calculate next state"})
	hooks = append(hooks, plugins.Hook{1, CoreSubStep1, nil, "core update next state"})

	return hooks
}

func (p *plugin) Filename() string {
	return pluginFilename
}

// note this logs through the context
func CoreSubStep1(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger).With("plugin", pluginFilename)
	occupiedCells := 0
	for _, cell := range life.cells {
		cell.occupied = cell.occupied_next
		if cell.occupied {
			occupiedCells += 1
		}
	}
	log.Info(fmt.Sprintf("%d alive cells", occupiedCells))
	return nil
}

// note this logs through the context
func ThreadSubStep0(ctx context.Context, id int, name string) error {
	//log := ctx.Value("log").(*slog.Logger).With("actor", name).With("plugin", pluginFilename)

	threads := 2 // this is hardcoded and the reason for the next big change

	// determine the chunksize
	chunkSize := (life.x * life.y) / threads
	if (life.x*life.y)%threads != 0 {
		chunkSize += 1
	}
	// iterative over this thread's chunk
	for idx := chunkSize * id; idx <= chunkSize*(id+1) && idx <= life.x*life.y; idx++ {

	}

	return nil
}
