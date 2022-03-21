package pinger

import (
   "testing"
   "time"
)

func Test(t *testing.T) {
   p, err := NewPinger("127.0.0.1", time.Millisecond * 900)
   if err != nil {
      t.Fatal(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         t.Error(err)
      }

      if !resp {
         t.Error("localhost didn't ping")
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
