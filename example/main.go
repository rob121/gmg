package main 


import (
    "log"
    "github.com/rob121/gmgsrv"
    "time"
    "flag"
    
)


var port string
var addr string

func main() {
    
    
    flag.StringVar(&addr,"address","","Grill Ip Address")
    flag.StringVar(&port,"port","8080","Grill Port")
    flag.Parse()
    
    grill := gmg.NewGrill(addr+":"+port,"5m") //5m debounce between event triggering ie on low pellet alarn, don't trigger event again for 5 minutes
    
    grill.GetId()
    
    grill.GetInfo()
    
    grill.Poll("5s") //duration
    
    grill.Event("grill.ready",func(g *gmg.Grill)(error){
        
        log.Println("Grill Ready")
        return nil
    })
    
    grill.Event("grill.main.temp.reached",func(g *gmg.Grill)(error){
        
        log.Println("Main Temp Reached")
        return nil
    })
    
    grill.Event("grill.probe1.temp.reached",func(g *gmg.Grill)(error){
        log.Println("Probe1 Temp Reached")
        return nil
    })
    
    grill.Event("grill.probe2.temp.reached",func(g *gmg.Grill)(error){
        log.Println("Probe2 Temp Reached")
        return nil
    })
    
    grill.Event("grill.pellets.low",func(g *gmg.Grill)(error){
        log.Println("Pellets Low")
        return nil
    })
    
    
    ticker := time.NewTicker(5 * time.Second)

    
    for _ = range ticker.C {
    log.Printf("%+v\n",grill)
    }
    
    
}