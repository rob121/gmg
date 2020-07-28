package gmg

import ( 
        "bufio"
     	"bytes"
        "net"
        "log"
        "fmt"
        "time"
        "errors"
        "encoding/hex"
        "github.com/AlexanderGrom/go-event"
)

const info ="UR001!"
const id = "UL!"

type GrillInfo struct {
     Connected bool
     Serial string
     State string
     On bool
     CurrentGrillTemp int64
     DesiredGrillTemp int64
     CurrentProbe1Temp int64
     DesiredProbe1Temp int64
     CurrentProbe2Temp int64
     DesiredProbe2Temp int64
     PelletsLow bool
}
    

type Grill struct {
    Info GrillInfo
    Settings struct {
       PizzaMode bool
    }
    oldinfo GrillInfo
    internal struct{
      address string
      event event.Dispatcher
      readyReached time.Time
      tempLastReached time.Time
      probe1LastReached time.Time
      probe2LastReached time.Time  
      tempEventDebounce time.Duration
      pelletsLowReached time.Time
      tempEventDebounceRegular time.Duration
    }
}

type GrillResponse struct{
    Bytes []byte
    Hex   string
}

func NewGrill(addr string,dur string) *Grill {
      
       
       
    g := new(Grill)
    
    if(dur==""){
        
        dur = "5m"
    }
    
    tdur,err := time.ParseDuration(dur)
    
    if(err!=nil){
        
        log.Fatal(err)
    }
    
    tdur2,_ := time.ParseDuration("4h")
    
    g.internal.address = addr
    g.internal.event = event.New()
    g.Info.Connected = false
    g.internal.tempEventDebounce = tdur //events we want to be reminded about regularly
    g.internal.tempEventDebounceRegular = tdur2 //events we only want to hear about 1x during cook
    g.internal.readyReached = time.Now().AddDate(0, -1, 0)
    g.internal.tempLastReached = time.Now().AddDate(0, -1, 0)
    g.internal.probe1LastReached = time.Now().AddDate(0, -1, 0)
    g.internal.probe2LastReached = time.Now().AddDate(0, -1, 0)
    g.internal.pelletsLowReached = time.Now().AddDate(0, -1, 0)
    return g
}    



func (g *Grill) GetId() {
    
   req,err := g.request(id)
   
   if(err!=nil){
       
       log.Println(err)
       
   }
   
   g.Info.Serial = string(bytes.Trim(req.Bytes, "\x00"))
   
   
}
 
func (g *Grill) Poll(durs string){
    
    dur,_ := time.ParseDuration(durs)
    
    ticker := time.NewTicker(dur)
	
    // for every `tick` that our `ticker`
    // emits, we print `tock`
    
    go func(){
    
	for _ = range ticker.C {
		g.GetInfo()
	}
	
	
	}()
    
    
}

func (g *Grill) Event(evt string,fn interface{}) (error){
    
return g.internal.event.On(evt, fn)
    
}

func (g *Grill) stateCheck(){
    
    
    if( (g.Info.CurrentGrillTemp >= g.Info.DesiredGrillTemp) && g.Info.CurrentGrillTemp > 150 && g.Info.DesiredGrillTemp > 155 && g.internal.tempLastReached.Add(g.internal.tempEventDebounceRegular).Before(time.Now()) ){
        
         g.internal.tempLastReached = time.Now()
         g.internal.event.Go("grill.main.temp.reached",g)

    }
    
    if( g.Info.CurrentGrillTemp >= 150 && g.Info.DesiredGrillTemp < 155 && g.internal.readyReached.Add(g.internal.tempEventDebounceRegular).Before(time.Now()) ){
        
         g.internal.readyReached = time.Now()
         g.internal.event.Go("grill.ready",g)

    }
    
    if(g.Info.CurrentProbe1Temp >= g.Info.DesiredProbe1Temp && g.Info.CurrentProbe1Temp > 0 && g.Info.DesiredProbe1Temp > 0 && g.internal.probe1LastReached.Add(g.internal.tempEventDebounce).Before(time.Now()) ){
        
        g.internal.probe1LastReached = time.Now()
        g.internal.event.Go("grill.probe1.temp.reached",g)
        
    }  
    
    if(g.Info.CurrentProbe2Temp >= g.Info.DesiredProbe2Temp && g.Info.CurrentProbe2Temp > 0 && g.Info.DesiredProbe2Temp > 0 && g.internal.probe2LastReached.Add(g.internal.tempEventDebounce).Before(time.Now()) ){
        
         g.internal.probe2LastReached = time.Now()
         g.internal.event.Go("grill.probe2.temp.reached",g)

        
    }  
    
    if(g.Info.PelletsLow == true && g.internal.pelletsLowReached.Add(g.internal.tempEventDebounce).Before(time.Now()) ){
        
         g.internal.pelletsLowReached = time.Now()
         g.internal.event.Go("grill.pellets.low",g)

    }
    
    if(g.oldinfo.State=="fan_mode" && g.Info.State=="off"){
        
        g.internal.event.Go("grill.cooldown.complete",g)  
    }
    
    
}

func (g *Grill) GetInfo() {
    
   if( g.Info!=GrillInfo{} ){
       
       g.oldinfo = g.Info
       
   } 
    
   req,err := g.request(info)
   
   if(err!=nil){
       
       log.Println(err)
       
   }
   
   if(len(req.Hex)>0){
   settings := req.Hex[17:32]
   
   g.Settings.PizzaMode = false
   
   if(settings[1:2]=="2"){
      
       g.Settings.PizzaMode = true 
       
   }
   }


  
 g.Info.State = getGrillState(req.Hex)
 g.Info.On = false
 if(g.Info.State=="on" || g.Info.State=="cold smoke mode"){
     
     g.Info.On = true
 }
 
 g.Info.CurrentGrillTemp  = getCurrentGrillTemp(req.Hex)
 g.Info.DesiredGrillTemp  = 0
 g.Info.CurrentProbe1Temp = 0
 g.Info.DesiredProbe1Temp = 0 
 g.Info.CurrentProbe2Temp = 0
 g.Info.DesiredProbe2Temp = 0 
 
 if(g.Info.State == "on"){
 
 g.Info.CurrentGrillTemp  = getCurrentGrillTemp(req.Hex)
 g.Info.DesiredGrillTemp= getDesiredGrillTemp(req.Hex)
 g.Info.CurrentProbe1Temp= getCurrentFoodTemp(req.Hex)
 g.Info.DesiredProbe1Temp= getDesiredFoodTemp(req.Hex) 
 g.Info.CurrentProbe2Temp= value(req.Hex,32).(int64)
 g.Info.DesiredProbe2Temp= value(req.Hex,36).(int64) 

 }


 g.Info.PelletsLow= getLowPelletAlarmActive(req.Hex)
 g.stateCheck()
    
}


func (g *Grill) request(cmd string) (GrillResponse, error){
    
    var buf bytes.Buffer
    gr := GrillResponse{}
	//fmt.Printf("%s    Request: Get All Info\n", time.Now().Format(time.RFC822))
	//fmt.Println("Request: Get All Info")
	fmt.Fprint(&buf, cmd)
	gResponse, err := g.rawrequest(&buf)
	
	if err != nil {
		return gr, err
	}
	
	
	gr.Bytes = gResponse
	gr.Hex = hex.EncodeToString(gResponse)
	
	return gr, nil
    
}


func (g *Grill) rawrequest(b *bytes.Buffer) ([]byte, error) {
	barray := make([]byte, 96)
	var err error
	var readBytes int
	retries := 5 // the grill doesnt always respond on the first try
	for i := 1; i <= retries; i++ {
		
		var conn net.Conn
		
		if i != 1 {
			fmt.Printf("Request Attempt %v\n", i)
		}
		
		if b.Len() == 0 && i == retries {
			return nil, errors.New("Nothing to Send to Grill")
		}
        
        conn, err = net.DialTimeout("udp", fmt.Sprintf("%s", g.internal.address), 3*time.Second)
		
		g.Info.Connected = false
		
		if err != nil {
		      	// TODO make this better
			fmt.Println("Error Connecting to Grill")
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			if i == retries {
				return nil, errors.New("1")
			}
			continue
		}

		
		if err != nil && i == retries {
			return nil, errors.New("Connection to Grill Failed")
		}
		
		g.Info.Connected = true
		
		timeout := time.Now().Add(time.Second)
		conn.SetReadDeadline(timeout) // sometimes the grill holds the conection forever
		//fmt.Println("Connected")

		defer conn.Close()
		//fmt.Println("Sending Data..")
		_, err = conn.Write(b.Bytes())
		if err != nil && i == retries {
                        fmt.Println(err)
			return nil, errors.New("Failure Sending Payload to Grill")
		}
		//fmt.Printf("Bytes Written: %v\n", ret)
		//b.Reset()

	   // fmt.Println("Reading Data..")
		readBytes, err = bufio.NewReader(conn).Read(barray)
		if err != nil && i == retries {
			return nil, errors.New("Failed Reading Result From Grill")
		}
		if readBytes > 0 {
			break
		}
	}
	barray = barray[:96]
	return barray, nil
}

