package pinger

import (
   "testing"
   "time"
)

func Test(t *testing.T) {
   p, err := NewPinger("localhost", time.Millisecond * 900)
   if err != nil {
      t.Fatal(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         t.Error(err)
      }

      if !resp {
         t.Fatal("no reply")
      }
   }

   p, err = NewPinger("128.0.0.1", time.Millisecond * 900)
   if err != nil {
      t.Fatal(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         t.Error(err)
      }

      if resp {
         t.Error("unexpected response")
      }
   }
}
