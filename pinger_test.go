package pinger

import (
   "testing"
   "time"
)

func Test(t *testing.T) {
   p, err := New("localhost", time.Millisecond * 900)
   if err != nil {
      t.Fatal(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         t.Error(err)
      }

      if resp == Timeout {
         t.Fatal("no reply")
      }
   }

   p, err = New("128.0.0.1", time.Millisecond * 900)
   if err != nil {
      t.Fatal(err)
   }

   for i := 0; i < 3; i++ {
      resp, err := p.Ping()
      if err != nil {
         t.Error(err)
      }

      if resp != Timeout {
         t.Error("unexpected response")
      }
   }
}
