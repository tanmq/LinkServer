package session

import (
  "sync"
  "errors"
  "strconv"

  "linkServer/logger"
  "linkServer/packet"
)

const (
  segSize = 16
)


type seg struct{
  sync.RWMutex
  m map[string]*Session
}

type SessionMap struct {
  segmap map[uint8]seg
}


func newSeg() seg {
  s := new(seg)
  s.m = make(map[string]*Session)

  return *s
}

func (s* seg) put(uid uint64, device packet.DeviceType, session *Session) {
  s.Lock()
  defer s.Unlock()

  key := hashItem(uid, device)
  _,ok := s.m[key]
  if ok {
    delete(s.m, key)
  }

  s.m[key] = session
}

func (s* seg)get(uid uint64, device packet.DeviceType) (*Session, error) {
  s.RLock()
  defer s.RUnlock()

  key := hashItem(uid, device)
  session, ok := s.m[key]
  if ok {
    return session, nil
  }

  return session, errors.New("Session does Exist")
}

func (s* seg) remove(uid uint64, device packet.DeviceType) {
  s.Lock()
  defer s.Unlock()

  key := hashItem(uid, device)

  delete(s.m, key)
}

func NewSessionMap() SessionMap {
  smap := SessionMap{}
  smap.segmap = make(map[uint8]seg)

  //init for each segment
  for i := 0; i < segSize; i ++ {
    s := newSeg()
    smap.segmap[uint8(i)] = s
  }

  return smap
}

func (s *SessionMap)Put(uid uint64, device packet.DeviceType, session *Session) error {
   skey := hashSeg(uid)
   seg, ok := s.segmap[skey]
   if ok {
     seg.put(uid, device, session)
     return nil
   }

   logger.Warn("Segment is not exist  un expected. skey=%d", skey)
   return errors.New("Insert failed. Segment not exist.")
}

func (s *SessionMap)Get(uid uint64, device packet.DeviceType) (*Session, error) {
  skey := hashSeg(uid)
  seg, ok := s.segmap[skey]
  if !ok {
    logger.Warn("Segment is not exist  un expected. skey=%d", skey)
    return nil, errors.New("uid not exist.")
  }

  return seg.get(uid, device)
}

func (s *SessionMap)Del(uid uint64, device packet.DeviceType) bool {
  skey := hashSeg(uid)
  seg, ok := s.segmap[skey]
  if ok {
    seg.remove(uid, device)
  }

  return true
}

func (s *SessionMap)SessionSize() int {
   count := 0
   for _,smap := range s.segmap {
     count += len(smap.m)
   }

   return count
}


//遍历所有session
func (s *SessionMap)Iter() chan<- *Session {
  iter := make(chan *Session)

  go func ()  {
    for _, v := range s.segmap {
      for _, session := range v.m {
        iter <- session
      }
    }

     close(iter)// finished iter
  }()

  return iter
}

//hash function for SESSIONMAP
func hashSeg(uid uint64) uint8 {
  return uint8(uid % segSize)
}

// hash fuction for SEGMAP
func hashItem(uid uint64, device packet.DeviceType) string {
  return strconv.FormatUint(uid, 10) + "-" + strconv.Itoa(int(device))
}
