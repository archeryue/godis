package main

import (
	"log"
	"time"

	"golang.org/x/sys/unix"
)

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
	id       int
	mask     TeType
	when     int64 //ms
	interval int64 //ms
	proc     TimeProc
	extra    interface{}
	next     *AeTimeEvent
}

type AeLoop struct {
	FileEvents      *AeFileEvent
	TimeEvents      *AeTimeEvent
	fdEventCnt      map[int]int
	fileEventFd     int
	timeEventNextId int
	stop            bool
}

func getEpollEvent(mask FeType) uint32 {
	if mask == AE_READABLE {
		return unix.EPOLLIN
	} else {
		return unix.EPOLLOUT
	}
}

func (loop *AeLoop) AddFileEvent(fd int, mask FeType, proc FileProc, extra interface{}) {
	// epoll ctl
	op := unix.EPOLL_CTL_ADD
	if loop.fdEventCnt[fd] > 0 {
		op = unix.EPOLL_CTL_MOD
	}
	err := unix.EpollCtl(loop.fileEventFd, op, fd, &unix.EpollEvent{Fd: int32(fd), Events: getEpollEvent(mask)})
	if err != nil {
		log.Printf("epoll ctr err: %v\n", err)
		return
	}
	loop.fdEventCnt[fd]++
	var fe AeFileEvent
	fe.fd = fd
	fe.mask = mask
	fe.proc = proc
	fe.extra = extra
	fe.next = loop.FileEvents
	loop.FileEvents = &fe
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
	err := unix.EpollCtl(loop.fileEventFd, unix.EPOLL_CTL_DEL, fd, &unix.EpollEvent{Fd: int32(fd), Events: getEpollEvent(mask)})
	if err != nil {
		log.Printf("epoll del err: %v\n", err)
	} else {
		loop.fdEventCnt[fd]--
	}
}

func GetMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func (loop *AeLoop) AddTimeEvent(mask TeType, interval int64, proc TimeProc, extra interface{}) int {
	id := loop.timeEventNextId
	loop.timeEventNextId++
	var te AeTimeEvent
	te.id = id
	te.mask = mask
	te.interval = interval
	te.when = GetMsTime() + interval
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

func AeLoopCreate() (*AeLoop, error) {
	epollFd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &AeLoop{
		fileEventFd:     epollFd,
		timeEventNextId: 1,
		fdEventCnt:      make(map[int]int),
		stop:            false,
	}, nil
}

func (loop *AeLoop) nearestTime() int64 {
	var nearest int64 = GetMsTime() + 1000
	p := loop.TimeEvents
	for p != nil {
		if p.when < nearest {
			nearest = p.when
		}
		p = p.next
	}
	return nearest
}

func (loop *AeLoop) getFileEvent(fd int, mask FeType) *AeFileEvent {
	p := loop.FileEvents
	for p != nil {
		if p.fd == fd && p.mask == mask {
			return p
		}
		p = p.next
	}
	return nil
}

func (loop *AeLoop) AeWait() (tes []*AeTimeEvent, fes []*AeFileEvent, err error) {
	timeout := loop.nearestTime() - GetMsTime()
	if timeout <= 0 {
		timeout = 10 // at least wait 10ms
	}
	var events [128]unix.EpollEvent
	n, err := unix.EpollWait(loop.fileEventFd, events[:], int(timeout))
	if err != nil {
		log.Printf("epoll wait err: %v\n", err)
		return
	}
	// collect file events
	for i := 0; i < n; i++ {
		if events[i].Events&unix.EPOLLIN != 0 {
			fe := loop.getFileEvent(int(events[i].Fd), AE_READABLE)
			if fe != nil {
				fes = append(fes, fe)
			}
		} else if events[i].Events&unix.EPOLLOUT != 0 {
			fe := loop.getFileEvent(int(events[i].Fd), AE_WRITABLE)
			if fe != nil {
				fes = append(fes, fe)
			}
		}
	}
	// collect time events
	now := GetMsTime()
	p := loop.TimeEvents
	for p != nil {
		if p.when <= now {
			tes = append(tes, p)
		}
		p = p.next
	}
	return
}

func (loop *AeLoop) AeProcess(tes []*AeTimeEvent, fes []*AeFileEvent) {
	for _, te := range tes {
		te.proc(loop, te.id, te.extra)
		if te.mask == AE_ONCE {
			loop.RemoveTimeEvent(te.id)
		} else {
			te.when = GetMsTime() + te.interval
		}
	}
	for _, fe := range fes {
		fe.proc(loop, fe.fd, fe.extra)
		loop.RemoveFileEvent(fe.fd, fe.mask)
	}
}

func (loop *AeLoop) AeMain() {
	for loop.stop != true {
		tes, fes, err := loop.AeWait()
		if err != nil {
			loop.stop = true
		}
		loop.AeProcess(tes, fes)
	}
}
