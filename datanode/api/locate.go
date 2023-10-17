package api

/*func Locate(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//todo 这块的工作代码当初设计是什么？
func StartLocate(mgfaceMQAddr, nodeAddr string) {
	MgfaceMQAddr = mgfaceMQAddr
	nodeAddr = nodeAddr
	client := NewTCPClient(mgfaceMQAddr)
	for {
		command := &Cmd{Name: "get", Key: "search", Value: ""}
		command.Run(client)
		if command.Value != "" {
			var mvalue []MetaValue
			json.Unmarshal([]byte(command.Value), &mvalue)
			for _, v := range mvalue {
				objectName := v.RealNodeValue
				if Locate(localStorePath + objectName) {
					command := &Cmd{Name: "set", Key: objectName, Value: nodeAddr}
					command.Run(client)
				}
			}
		}
		time.Sleep(time.Second)
	}
}*/
