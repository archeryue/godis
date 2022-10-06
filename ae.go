package main

type FeType int

const (
	AE_READABLE FeType = 1
	AE_WRITABLE FeType = 2
)

type TeType int

const (
	AE_NORMAL TeType = 1
	AE_ONCE   TeType = 2
)

type FileProc func(loop *AeLoop, fd int, extra interface{})
type TimeProc func(loop *AeLoop, id int, extra interface{})

type AeFileEvent struct {
	fd    int
	mask  FeType
	proc  FileProc
	extra interface{}
	next  *AeFileEvent
}

type AeTimeEvent struct {
	id    int
	mask  TeType
	when  int64 //second
	proc  TimeProc
	extra interface{}
	next  *AeTimeEvent
}

type AeLoop struct {
	FileEvents      *AeFileEvent
	TimeEvents      *AeTimeEvent
	timeEventNextId int
	stop            bool
}

func (loop *AeLoop) AddFileEvent(fd int, mask FeType, proc FileProc, extra interface{}) {
	var fe AeFileEvent
	fe.fd = fd
	fe.mask = mask
	fe.proc = proc
	fe.extra = extra
	fe.next = loop.FileEvents
	loop.FileEvents = &fe
	//TODO: epoll ctl
}

func (loop *AeLoop) RemoveFileEvent(fd int, mask FeType) {
	p := loop.FileEvents
	var pre *AeFileEvent
	for p != nil {
		if p.fd == fd && p.mask == mask {
			if pre == nil {
				loop.FileEvents = p.next
			} else {
				pre.next = p.next
			}
			p.next = nil
			break
		}
		pre = p
		p = p.next
	}
	// TODO: epoll ctl
}

func (loop *AeLoop) AddTimeEvent(mask TeType, when int64, proc TimeProc, extra interface{}) int {
	id := loop.timeEventNextId
	loop.timeEventNextId++
	var te AeTimeEvent
	te.id = id
	te.mask = mask
	te.when = when
	te.proc = proc
	te.extra = extra
	te.next = loop.TimeEvents
	loop.TimeEvents = &te
	return id
}

func (loop *AeLoop) RemoveTimeEvent(id int) {
	p := loop.TimeEvents
	var pre *AeTimeEvent
	for p != nil {
		if p.id == id {
			if pre == nil {
				loop.TimeEvents = p.next
			} else {
				pre.next = p.next
			}
			p.next = nil
			break
		}
		pre = p
		p = p.next
	}
}

func AeLoopCreate() *AeLoop {
	var loop AeLoop
	loop.timeEventNextId = 1
	loop.stop = false
	return &loop
}

func (loop *AeLoop) AeWait() (tes []AeTimeEvent, fes []AeFileEvent) {
	//TODO: search time && epoll wait
	return nil, nil
}

func (loop *AeLoop) AeProcess(tes []AeTimeEvent, fes []AeFileEvent) {
	for _, te := range tes {
		te.proc(loop, te.id, te.extra)
		if te.mask == AE_ONCE {
			loop.RemoveTimeEvent(te.id)
		}
	}
	for _, fe := range fes {
		fe.proc(loop, fe.fd, fe.extra)
		loop.RemoveFileEvent(fe.fd, fe.mask)
	}
}

func (loop *AeLoop) AeMain() {
	for loop.stop != true {
		tes, fes := loop.AeWait()
		loop.AeProcess(tes, fes)
	}
}
