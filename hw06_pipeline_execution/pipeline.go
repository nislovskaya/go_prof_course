package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) Out

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = mergeChan(done, stage(out))
	}
	return out
}

func mergeChan(done, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			case <-done:
				return
			}
		}
	}()

	return out
}
