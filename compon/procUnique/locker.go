package procUnique

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

type Locker struct {
	name     string
	filePath string
	file     *os.File
}

func NewLocker(inname string) *Locker {
	lock := new(Locker)

	tmpDir := os.TempDir()
	lockFileName := fmt.Sprintf("%s_uniqProcLock.pid", inname)
	filePath := path.Join(tmpDir, lockFileName)
	lock.filePath = filePath
	lock.name = inname
	lock.file = nil
	return lock
}

func (this *Locker) Unlock() {
	if this.file != nil {
		syscall.Flock(int(this.file.Fd()), syscall.LOCK_UN)
		this.file.Close()

		if len(this.filePath) > 0 {
			os.Remove(this.filePath)
		}
	}
}

func (this *Locker) Lock() error {
	lockFile, err := os.Create(this.filePath)
	if err != nil {
		this.file = nil
		return err
	}

	//syscall.LOCK_EX 排它锁，不允许其他人读和写。syscall.LOCK_NB 意味着无法锁定文件时不能阻断操作，马上返回给进程。
	//lock.Fd()返回文件描述符，文件描述符是一个索引值，指向当前进程打开的文件记录表。最后在执行完毕后对文件解锁。
	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		lockFile.Close()
		this.file = nil
		return err
	}

	lockFile.WriteString(fmt.Sprint(os.Getpid()))
	this.file = lockFile
	return nil
}
