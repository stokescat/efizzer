package efi

import (
  "fmt"
)

type Guid string

var (
  NotFvGuid       = Guid("500DF8D1-CA05-3042-959E-8872D94F6BC6")
  UndefinedGuid   = Guid("336F5A24-E9F9-B643-972A-14C8738E3F5B") 
)

func BytesToGuid(b []byte) Guid {
  
  guidStr:= fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
      b[3],b[2],b[1],b[0], b[5],b[4], b[7],b[6], b[8],b[9], b[10],b[11],b[12],b[13],b[14],b[15])
  
  return Guid(guidStr) 
}
