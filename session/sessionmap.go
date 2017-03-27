package session

import (
	"errors"
	"strconv"
	"sync"

	"linkServer/logger"
	"linkServer/packet"
)

const (
	segSize = 16 //cut sessionmap into 16 segment for low lock level.
)

type segment struct {
	sync.RWMutex
	m map[string]*Session
}

//SSPool store all user sessions.
type SSPool struct {
	segmap map[uint8]*segment
}

func newSeg() *segment {
	s := new(segment)
	s.m = make(map[string]*Session)

	return s
}

func (s *segment) put(uid uint64, device packet.DeviceType, ss *Session) {
	s.Lock()
	defer s.Unlock()

	key := buildSegmentKey(uid, device)
	s.closeSession(key) // close old one if exist

	s.m[key] = ss
}

func (s *segment) size() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.m)
}

func (s *segment) get(uid uint64, device packet.DeviceType) (*Session, error) {
	s.RLock()
	defer s.RUnlock()

	key := buildSegmentKey(uid, device)
	ss, ok := s.m[key]
	if ok {
		return ss, nil
	}

	return ss, errors.New("Session does Exist")
}

//should be locked s.m by who calls
func (s *segment) closeSession(key string) {
	ss, ok := s.m[key]
	if ok {
		delete(s.m, key)
		ss.Close() // close tcp connection.
	}
}

func (s *segment) remove(uid uint64, device packet.DeviceType) {
	s.Lock()
	defer s.Unlock()

	key := buildSegmentKey(uid, device)
	s.closeSession(key)
}

//NewSSPool construct new session pool to store user sessions.
func NewSSPool() SSPool {
	smap := SSPool{}
	smap.segmap = make(map[uint8]*segment)

	//init for each segment
	for i := 0; i < segSize; i++ {
		s := newSeg()
		smap.segmap[uint8(i)] = s
	}

	return smap
}

// Put add a new session into pool.
func (s *SSPool) Put(uid uint64, device packet.DeviceType, ss *Session) error {
	skey := buildSessionMapKey(uid)
	seg, ok := s.segmap[skey]
	if ok {
		seg.put(uid, device, ss)
		return nil
	}

	logger.Warn("Segment is not exist  un expected. skey=%d", skey)
	return errors.New("segment not exist")
}

//Get fetch a user session.
func (s *SSPool) Get(uid uint64, device packet.DeviceType) (*Session, error) {
	skey := buildSessionMapKey(uid)
	seg, ok := s.segmap[skey]
	if !ok {
		logger.Warn("Segment is not exist  un expected. skey=%d", skey)
		return nil, errors.New("uid not exist")
	}

	return seg.get(uid, device)
}

//Del delete a user session and close it.
func (s *SSPool) Del(uid uint64, device packet.DeviceType) bool {
	skey := buildSessionMapKey(uid)
	seg, ok := s.segmap[skey]
	if ok {
		seg.remove(uid, device)
	}

	return true
}

//PoolSize return the number of sessions
func (s *SSPool) PoolSize() int {
	count := 0
	for _, smap := range s.segmap {
		count += smap.size()
	}

	return count
}

//Iter will return all sessions in turns by channel.
func (s *SSPool) Iter() chan<- *Session {
	iter := make(chan *Session)

	go func() {
		for _, v := range s.segmap {
			v.RLock() // should protect the segmap
			for _, ss := range v.m {
				iter <- ss
			}
			v.RUnlock()
		}

		close(iter) // finished iter
	}()

	return iter
}

func buildSessionMapKey(uid uint64) uint8 {
	return uint8(uid % segSize)
}

func buildSegmentKey(uid uint64, device packet.DeviceType) string {
	return strconv.FormatUint(uid, 10) + "-" + strconv.Itoa(int(device))
}
