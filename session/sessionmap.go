package session

import (
  "sync"
  "errors"
  "strconv"

  "linkServer/logger"
  "linkServer/packet"
)

const (
  segSize = 16  //cut sessionmap into 16 segment for low lock level.
)


type segment struct{
  sync.RWMutex
  m map[string]*Session
}

type SessionMap struct {
  segmap map[uint8]*segment
}


func newSeg() *segment {
  s := new(segment)
  s.m = make(map[string]*Session)

  return s
}

func (s* segment) put(uid uint64, device packet.DeviceType, ss *Session) {
  s.Lock()
  defer s.Unlock()

  key := buildSegmentKey(uid, device)
  s.closeSession(key) // close old one if exist

  s.m[key] = ss
}

func (s *segment)size() int{
  s.RLock()
  defer s.RUnlock()

  return len(s.m)
}

func (s* segment)get(uid uint64, device packet.DeviceType) (*Session, error) {
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
func (s *segment)closeSession(key string) {
  ss, ok := s.m[key]
  if ok {
    delete(s.m, key)
    ss.Close() // close tcp connection.
  }
}

func (s* segment) remove(uid uint64, device packet.DeviceType) {
  s.Lock()
  defer s.Unlock()

  key := buildSegmentKey(uid, device)
  s.closeSession(key)
}

func NewSessionMap() SessionMap {
  smap := SessionMap{}
  smap.segmap = make(map[uint8]*segment)

  //init for each segment
  for i := 0; i < segSize; i ++ {
    s := newSeg()
    smap.segmap[uint8(i)] = s
  }

  return smap
}

func (s *SessionMap)Put(uid uint64, device packet.DeviceType, ss *Session) error {
   skey := buildSessionMapKey(uid)
   seg, ok := s.segmap[skey]
   if ok {
     seg.put(uid, device, ss)
     return nil
   }

   logger.Warn("Segment is not exist  un expected. skey=%d", skey)
   return errors.New("Insert failed. Segment not exist.")
}

func (s *SessionMap)Get(uid uint64, device packet.DeviceType) (*Session, error) {
  skey := buildSessionMapKey(uid)
  seg, ok := s.segmap[skey]
  if !ok {
    logger.Warn("Segment is not exist  un expected. skey=%d", skey)
    return nil, errors.New("uid not exist.")
  }

  return seg.get(uid, device)
}

func (s *SessionMap)Del(uid uint64, device packet.DeviceType) bool {
  skey := buildSessionMapKey(uid)
  seg, ok := s.segmap[skey]
  if ok {
    seg.remove(uid, device)
  }

  return true
}

func (s *SessionMap)SessionSize() int {
   count := 0
   for _,smap := range s.segmap {
     count += smap.size()
   }

   return count
}


//遍历所有session
//这里会有个隐患，由于遍历需要加锁，如果处理不及时会导致加锁时间过长
func (s *SessionMap)Iter() chan<- *Session {
  iter := make(chan *Session)

  go func ()  {
    for _, v := range s.segmap {
      v.RLock() // should protect the segmap
      for _, ss := range v.m {
        iter <- ss
      }
      v.RUnlock()
    }

     close(iter)// finished iter
  }()

  return iter
}

func buildSessionMapKey(uid uint64) uint8 {
  return uint8(uid % segSize)
}

func buildSegmentKey(uid uint64, device packet.DeviceType) string {
  return strconv.FormatUint(uid, 10) + "-" + strconv.Itoa(int(device))
}
