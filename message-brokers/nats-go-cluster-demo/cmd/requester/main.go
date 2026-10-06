package main
import("fmt";"log";"os";"strings";"time";"github.com/nats-io/nats.go")
func main(){
 urls:=env("NATS_URLS","nats://localhost:4222,nats://localhost:4223,nats://localhost:4224")
 nc,err:=nats.Connect(strings.ReplaceAll(urls,","," "),nats.Name("requester"),nats.MaxReconnects(-1))
 if err!=nil{log.Fatal(err)}; defer nc.Close()
 for i:=1;i<=20;i++{
  req:=fmt.Sprintf("request-%d",i)
  msg,err:=nc.Request("math.double",[]byte(req),2*time.Second)
  if err!=nil{log.Printf("request failed: %v",err)}else{log.Printf("%s -> %s",req,string(msg.Data))}
  time.Sleep(time.Second)
 }
}
func env(k,v string)string{if x:=os.Getenv(k);x!=""{return x};return v}
