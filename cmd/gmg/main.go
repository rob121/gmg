package main 


import (
    "log"
    "github.com/rob121/gmg"
    "github.com/rob121/vhelp"
    "time"
    "flag"
    "github.com/davecgh/go-spew/spew"
    "github.com/spf13/viper"
    "net/http"
)

var grill *gmg.Grill
var port string
var conf *viper.Viper
var page_port string
var addr string
var assets_dir string

func main() {

    vhelp.Load("config")
    var err error

    conf,err = vhelp.Get("config")

    conf.SetDefault("page_port","8712")

    if(err!=nil){

        log.Fatal(err)
    }

    flag.StringVar(&addr,"address","","Grill Ip Address")
    flag.StringVar(&port,"port","8080","Grill Port")
    flag.Parse()

    if(len(addr)<1){

        log.Fatal("Address required")
    }

    go startHTTPServer()
    
    grill = gmg.NewGrill(addr+":"+port,"5m")
    
    grill.GetId()
    
    grill.GetInfo()
    
    grill.Poll("5s") //duration
    
    grill.Event("grill.ready",func(g *gmg.Grill)(error){
        
        log.Println("Grill Ready")

        url := conf.GetString("event.grill_ready")
        http.Get(url)
        return nil
    })
    
    grill.Event("grill.main.temp.reached",func(g *gmg.Grill)(error){

        url := conf.GetString("event.main_temp_reached")
        http.Get(url)

        log.Println("Main Temp Reached")
        return nil
    })
    
    grill.Event("grill.probe1.temp.reached",func(g *gmg.Grill)(error){

        url := conf.GetString("event.probe1_temp_reached")
        http.Get(url)

        log.Println("Probe1 Temp Reached")
        return nil
    })
    
    grill.Event("grill.probe2.temp.reached",func(g *gmg.Grill)(error){

        url := conf.GetString("event.probe2_temp_reached")
        http.Get(url)

        log.Println("Probe2 Temp Reached")
        return nil
    })
    
    grill.Event("grill.pellets.low",func(g *gmg.Grill)(error){

        url := conf.GetString("event.pellets_low")
        http.Get(url)

        log.Println("Pellets Low")

        
        return nil
        
    })
    
    grill.Event("grill.cooldown.complete",func(g *gmg.Grill)(error){

        url := conf.GetString("event.cooldown_complete")
        http.Get(url)

        log.Println("Pellets Low")

        
        return nil
        
    })
    
    
    ticker := time.NewTicker(5 * time.Second)

    
    for _ = range ticker.C {
    spew.Dump(grill.Info)
    }


    
    
}