package main

import (
  "flag"
  "fmt"
  "encoding/json"
  "github.com/cloudfoundry/gosigar"
  "time"
  "net/http"
  "bytes"
  "os"
)

func main() {

  var dashboard_hostname = flag.String("dashboard-hostname", "dashboard.comparaonline.com", "dashing server hostname")
  var auth_token = flag.String("auth-token", "", "dashing server auth token")
  flag.Parse()
  hostname,_ := os.Hostname()

  for { 
    avg := sigar.LoadAverage{}
    avg.Get()

    mem := sigar.Mem{}
    mem.Get()
    memUsagePercent := float64(mem.ActualUsed)/float64(mem.Total)*100

    maxHddUsage := 0.0
    fslist := sigar.FileSystemList{}
    fslist.Get()
    for _, fs := range fslist.List {
      dir_name := fs.DirName
      usage := sigar.FileSystemUsage{}
      usage.Get(dir_name)
      usePercent := usage.UsePercent()
      if(usePercent > maxHddUsage) {
        maxHddUsage = usePercent
      }
    }
    
    data := map[string]string{
      "hostname": hostname,
      "load": fmt.Sprintf("%.2f", avg.One),
      "mem": fmt.Sprintf("%.1f", memUsagePercent),
      "hdd": fmt.Sprintf("%.1f", maxHddUsage),
      "auth_token": *auth_token,
    }
    b,_ := json.Marshal(data)

    req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/widgets/server-%s", *dashboard_hostname, hostname), bytes.NewReader(b))
    client := &http.Client{}
    client.Do(req)

    fmt.Println(time.Now(),data)

    time.Sleep(15*time.Second)
  }
}
