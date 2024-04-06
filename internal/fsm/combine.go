package fsm

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func NewCombine(resultMachine Machine, machines ...Machine) *Combine {
	return &Combine{
		machines:      machines,
		resultMachine: resultMachine,
	}
}

type Combine struct {
	machines      []Machine
	resultMachine Machine
}

func (c *Combine) Switch(ctx context.Context, state State) (context.Context, Machine, State, error) {
	g, ctx := errgroup.WithContext(ctx)

	for _, machine := range c.machines {
		g.Go(func() error {
			st := state
			var err error
			for machine != nil {
				ctx, machine, st, err = machine.Switch(ctx, st)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return ctx, nil, state, err
	}
	//TODO: How combine state from all machines?

	return ctx, c.resultMachine, state, nil
}
