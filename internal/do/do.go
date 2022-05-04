package do

// History represents a stack of actions that can be undone or redone.
type History[T any] struct {
	Target   T           // Target is the underlying data that should be changed through the application or reverse application of actions.
	Stack    []Action[T] // Stack is the current stack of actions.
	Pos      int         // Pos is the internal action position within the stack.
	SavedPos int         // SavedPos is the last saved position, used externally.
}

// Reset empties the history and sets the position to 0.
func (d *History[T]) Reset() {
	d.Stack = make([]Action[T], 0)
	d.Pos = 0
	d.SavedPos = 0
}

// Push pushes a new action onto the stack.
func (d *History[T]) Push(a Action[T]) {
	d.Stack = append(d.Stack[:d.Pos], a)
	if d.SavedPos > d.Pos {
		d.SavedPos = -1
	}
	d.Pos++
}

// PushAndApply pushes a new action onto the stack and calls its Apply method.
func (d *History[T]) PushAndApply(a Action[T]) {
	d.Push(a)
	d.Stack[len(d.Stack)-1].Apply(d.Target)
}

// Undo unapplies the current action and decrements the position if possible.
func (d *History[T]) Undo() {
	if d.Pos-1 >= 0 {
		d.Stack[d.Pos-1].Unapply(d.Target)
		d.Pos--
	}
}

// Redo applies the next stored action and increments the position if one exists.
func (d *History[T]) Redo() {
	if d.Pos < len(d.Stack) {
		d.Stack[d.Pos].Apply(d.Target)
		d.Pos++
	}
}

// Undoable returns if Undo is able to be called.
func (d *History[T]) Undoable() bool {
	return d.Pos > 0 && len(d.Stack) > 0
}

// Redoable returns if Redo is able to be called.
func (d *History[T]) Redoable() bool {
	return d.Pos < len(d.Stack)
}

// Action is any redoable or undoable state. Apply should change the underlying Target to conform to whatever data changes the Action represents. Unapply should cleanly reverse the result of Apply.
type Action[T any] interface {
	Unapply(target T)
	Apply(target T)
}
