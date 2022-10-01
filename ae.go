package main

type FileProc func(loop AeLoop, fd int, mask int)
type TimeProc func(loop AeLoop, id int) int

type AeFileEvent struct {
	fd   int
	mask int
	proc FileProc
	next *AeFileEvent
}

type AeTimeEvent struct {
	id   int
	when int64
	proc TimeProc
	next *AeTimeEvent
}

type AeLoop struct {
	FileEvents      *AeFileEvent
	TimeEvents      *AeTimeEvent
	timeEventNextId int
	stop            int
}
