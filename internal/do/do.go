package do

type History[T any] struct {
	Target T
	Stack  []Action[T]
	Pos    int
}

func (d *History[T]) Push(a Action[T]) {
	d.Stack = append(d.Stack[:d.Pos], a)
	d.Pos++
}

func (d *History[T]) PushAndApply(a Action[T]) {
	d.Push(a)
	d.Stack[len(d.Stack)-1].Apply(d.Target)
}

func (d *History[T]) Undo() {
	if d.Pos-1 >= 0 {
		d.Stack[d.Pos-1].Unapply(d.Target)
		d.Pos--
	}
}

func (d *History[T]) Redo() {
	if d.Pos < len(d.Stack) {
		d.Stack[d.Pos].Apply(d.Target)
		d.Pos++
	}
}

func (d *History[T]) Undoable() bool {
	return d.Pos > 0 && len(d.Stack) > 0
}

func (d *History[T]) Redoable() bool {
	return d.Pos < len(d.Stack)
}

type Action[T any] interface {
	Unapply(target T)
	Apply(target T)
}
