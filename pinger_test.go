package pinger

import (
   "testing"
   "time"
)

func Test(t *testing.T) {
   p, err := NewPinger("127.0.0.1", time.Millisecond * 900)
   if err != nil {
      panic(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         panic(err)
      }

      if !resp {
         t.Fatal("localhost didn't ping")
      }
   }

   p, err = NewPinger("128.0.0.1", time.Millisecond * 900)
   if err != nil {
      panic(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         panic(err)
      }

      if resp {
         t.Fatal("unexpected response")
      }
   }
}
