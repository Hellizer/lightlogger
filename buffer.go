package logger

type LogObject struct {
	LogMsg
	next *LogObject
}

type linkedQueue struct {
	current *LogObject
	last    *LogObject
	cap     uint64
}

func (lq *linkedQueue) put(o *LogObject) {
	if lq.cap == 0 {
		lq.current = o
		lq.last = o
		lq.cap = 1
		return
	}
	lq.last.next = o
	lq.last = o
	lq.cap++
	return
}

func (lq *linkedQueue) get() *LogObject {
	if lq.cap == 0 {
		return nil
	}
	o := *lq.current
	lq.current = lq.current.next
	lq.cap--
	return &o
}
