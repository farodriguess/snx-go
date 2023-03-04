package snx

import (
	"fmt"
	"math/big"
	"net"
	"os/exec"
	"strconv"

	binary_pack "github.com/roman-kachanovsky/go-binary-pack/binary-pack"
)

type SNXExtender struct {
	Params map[string]string
	Debug  bool
	info   []byte
}

func (extender *SNXExtender) CallSNX() {
	extender.generateSNXInfo()
	extender.callSNX()
}

func (extender *SNXExtender) generateSNXInfo() {

	extender.log("generateSNXInfo start")
	params := extender.Params

	gwIP, err := net.LookupHost(params["host_name"])
	checkError(err)

	bp := new(binary_pack.BinaryPack)

	ip := net.ParseIP(gwIP[0])
	ipv4 := big.NewInt(0)
	ipv4.SetBytes(ip.To4())
	tmp := ip.To4()

	hwData, err := bp.UnPack([]string{"I"}, []byte{tmp[3], tmp[2], tmp[1], tmp[0]})
	checkError(err)

	gwInt := hwData[0].(int)

	magic := string([]byte{0x13, 0x11, 0x00, 0x00})
	length := 0x3d0

	port, err := strconv.Atoi(params["port"])
	checkError(err)

	format := []string{"4s", "L", "L", "64s", "L", "6s", "256s", "256s", "128s", "256s", "H"}

	values := []interface{}{
		magic,
		length,
		gwInt,
		params["host_name"],
		port,
		string([]byte{0}),
		params["server_cn"],
		params["user_name"],
		params["password"],
		params["server_fingerprint"],
		1,
	}

	data, err := bp.Pack(format, values)
	checkError(err)

	extender.log("Packed Values: %v", data)
	extender.log("generateSNXInfo end")
	extender.info = data

}

func (extender *SNXExtender) callSNX() {

	extender.log("callSNX start")

	snxCmd := exec.Command("/usr/bin/snx", "-Z")

	_, err := snxCmd.Output()
	checkError(err)

	connection, err := net.Dial("tcp", "localhost:7776")
	checkError(err)

	extender.log("writing info %v", extender.info)
	_, err = connection.Write(extender.info)
	checkError(err)

	buffer := make([]byte, 4096)

	_, err = connection.Read(buffer)
	checkError(err)

	extender.log("snxanswer: %v", buffer)

	fmt.Println("SNX connected, to leave VPN open, leave this running!")
	connection.Read(buffer) //Block execution

	extender.log("callSNX end")

}

func (extender *SNXExtender) log(msg string, a ...any) {
	if extender.Debug {
		fmt.Println(fmt.Sprintf(msg, a...))
	}
}
