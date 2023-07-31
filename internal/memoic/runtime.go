package memoic

// Runtime contains the current Runtime information such as the parameters,
// the heap, not to be confused by Stack's sector and the end result of the
// Runtime. Unlike Stack, the heap is protected and can only be written by
// the runtime.
type Runtime struct {
	Parameters Sector
	heap       Sector
	stacks     []Stack
	Result     *any
}

// Get retrieves an item from the heap.
func (runtime *Runtime) Get(key string) any {
	if value, ok := runtime.heap[key]; ok {
		return value
	}
	return nil
}

// Load loads a function with the given parameters.
func (runtime *Runtime) Load(fn *Function) error {
	for _, pipe := range *fn {
		stack := runtime.newStack()
		stack.value = pipe.Value
		stack.interpolate()
		result, err := pipe.Invoke(stack)
		if err != nil {
			return err
		}
		if result != nil && pipe.As != nil {
			runtime.heap[*pipe.As] = result
		}
	}
	return nil
}

func (runtime *Runtime) newStack() *Stack {
	stack := Stack{Runtime: runtime}
	runtime.stacks = append(runtime.stacks, stack)
	return &stack
}

func NewRuntime(parameters Sector) *Runtime {
	return &Runtime{Parameters: parameters, heap: Sector{}}
}
