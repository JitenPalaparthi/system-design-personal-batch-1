package main
import("flag";"fmt";"log";"os";"strings";"time";"github.com/nats-io/nats.go")
func main(){
 var name string; flag.StringVar(&name,"name","service-1","service name");flag.Parse()
 urls:=env("NATS_URLS","nats://localhost:4222,nats://localhost:4223,nats://localhost:4224")
 nc,err:=nats.Connect(strings.ReplaceAll(urls,","," "),nats.Name(name),nats.MaxReconnects(-1),nats.ReconnectWait(time.Second))
 if err!=nil{log.Fatal(err)};defer nc.Drain()
 _,err=nc.QueueSubscribe("math.double","math-services",func(m *nats.Msg){
  log.Printf("[%s] request=%s reply=%s",name,string(m.Data),m.Reply)
  _=m.Respond([]byte(fmt.Sprintf("handled-by=%s; result-for=%s",name,string(m.Data))))
 })
 if err!=nil{log.Fatal(err)};nc.Flush()
 log.Printf("[%s] listening on math.double queue=math-services",name);select{}
}
func env(k,v string)string{if x:=os.Getenv(k);x!=""{return x};return v}
