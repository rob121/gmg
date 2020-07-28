package gmg 

import (
     "strconv"   
     "log"

)

func getRawValue (hexstr string, start int,end int) int64 {
  

    
 dig,_:=strconv.ParseInt(hexstr[start:end],16,64)
 return dig
}

func getCurrentGrillTemp (hexstr string) int64 {
   if(len(hexstr)<1){
       return 0
   } 
    
  first := getRawValue(hexstr, 4,6);
  second := getRawValue(hexstr, 6,8);
  return first + second * 256;
}

func getLowPelletAlarmActive(hexstr string) bool {
   if(len(hexstr)<1){
       return false
   }  
    
   first := getRawValue(hexstr, 48,50);
   second := getRawValue(hexstr, 50,52);
   value := first + second * 256;
   if( value == 128 ){
       return true
   }
   return false
}


func getDesiredGrillTemp(hexstr string) int64{
   if(len(hexstr)<1){
       return 0
   }  
    
   first := getRawValue(hexstr, 12,14);
   second := getRawValue(hexstr, 14,16);
   return first + second * 256;
}

func getGrillState (hexstr string) string {
  if(len(hexstr)<1){
       return "powered_down"
  } 
  
  status,_ := strconv.ParseInt(string(hexstr[61]),10,64)
  statusr := ""
  if (status == 0){
    statusr = "off"
  }else if(status == 1){ 
       statusr = "on" 
       
  }else if(status == 2){ 
      statusr = "fan_mode" 
      
  }else if(status == 3){ 
      statusr = "cold_smoke_mode" 
  }else{ 
      statusr = "unknown" 
  }
  
  log.Println("Status:",status)
  
  return statusr
}

func getCurrentFoodTemp (hexstr string) int64  {
   if(len(hexstr)<1){
       return 0
   }  
    
   first := getRawValue(hexstr, 8,10)
   second := getRawValue(hexstr, 10,12)
   currentFoodTemp := first + second * 256
    if( currentFoodTemp >= 557){
    return 0 
    }else{
    return currentFoodTemp
    }
}

func getDesiredFoodTemp (hexstr string)  int64 {
   if(len(hexstr)<1){
       return 0
   }  
    
   first := getRawValue(hexstr, 56,58)
   second := getRawValue(hexstr, 58,60)
  return first + second * 256
}

func value (hexstr string, start int) interface{}{
   if(len(hexstr)<1){
       return nil
   }  
    
   first := getRawValue(hexstr, start,start+2)
   second := getRawValue(hexstr, start + 2,start+4)
   return first + second * 256
}