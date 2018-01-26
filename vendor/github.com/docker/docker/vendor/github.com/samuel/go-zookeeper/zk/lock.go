package zk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrDeadlock  = errors.New("zk: trying to acquire a lock twice")
	ErrNotLocked = errors.New("zk: not locked")
)

type Lock struct ***REMOVED***
	c        *Conn
	path     string
	acl      []ACL
	lockPath string
	seq      int
***REMOVED***

func NewLock(c *Conn, path string, acl []ACL) *Lock ***REMOVED***
	return &Lock***REMOVED***
		c:    c,
		path: path,
		acl:  acl,
	***REMOVED***
***REMOVED***

func parseSeq(path string) (int, error) ***REMOVED***
	parts := strings.Split(path, "-")
	return strconv.Atoi(parts[len(parts)-1])
***REMOVED***

func (l *Lock) Lock() error ***REMOVED***
	if l.lockPath != "" ***REMOVED***
		return ErrDeadlock
	***REMOVED***

	prefix := fmt.Sprintf("%s/lock-", l.path)

	path := ""
	var err error
	for i := 0; i < 3; i++ ***REMOVED***
		path, err = l.c.CreateProtectedEphemeralSequential(prefix, []byte***REMOVED******REMOVED***, l.acl)
		if err == ErrNoNode ***REMOVED***
			// Create parent node.
			parts := strings.Split(l.path, "/")
			pth := ""
			for _, p := range parts[1:] ***REMOVED***
				pth += "/" + p
				_, err := l.c.Create(pth, []byte***REMOVED******REMOVED***, 0, l.acl)
				if err != nil && err != ErrNodeExists ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if err == nil ***REMOVED***
			break
		***REMOVED*** else ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	seq, err := parseSeq(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for ***REMOVED***
		children, _, err := l.c.Children(l.path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		lowestSeq := seq
		prevSeq := 0
		prevSeqPath := ""
		for _, p := range children ***REMOVED***
			s, err := parseSeq(p)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if s < lowestSeq ***REMOVED***
				lowestSeq = s
			***REMOVED***
			if s < seq && s > prevSeq ***REMOVED***
				prevSeq = s
				prevSeqPath = p
			***REMOVED***
		***REMOVED***

		if seq == lowestSeq ***REMOVED***
			// Acquired the lock
			break
		***REMOVED***

		// Wait on the node next in line for the lock
		_, _, ch, err := l.c.GetW(l.path + "/" + prevSeqPath)
		if err != nil && err != ErrNoNode ***REMOVED***
			return err
		***REMOVED*** else if err != nil && err == ErrNoNode ***REMOVED***
			// try again
			continue
		***REMOVED***

		ev := <-ch
		if ev.Err != nil ***REMOVED***
			return ev.Err
		***REMOVED***
	***REMOVED***

	l.seq = seq
	l.lockPath = path
	return nil
***REMOVED***

func (l *Lock) Unlock() error ***REMOVED***
	if l.lockPath == "" ***REMOVED***
		return ErrNotLocked
	***REMOVED***
	if err := l.c.Delete(l.lockPath, -1); err != nil ***REMOVED***
		return err
	***REMOVED***
	l.lockPath = ""
	l.seq = 0
	return nil
***REMOVED***
