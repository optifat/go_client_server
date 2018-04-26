//https://habr.com/post/126461/

package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha1"
	"strconv"
	"math/big"
)

const (
	tcpProtocol = "tcp4"
	keySize = 1024
	readWriterSize = keySize/8
)

type remoteConn struct {
	c *net.TCPConn
	pubK *rsa.PublicKey
}


func checkErr(err error){
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

var listenAddr = &net.TCPAddr{IP: net.IPv4(192, 168, 0, 108), Port: 0}

func waitPubKey(buf *bufio.Reader)(*rsa.PublicKey){
	line, _ , err := buf.ReadLine(); checkErr(err);
	if string(line) == "Connect"{
		line, _, err = buf.ReadLine(); checkErr(err)
		pubKey := rsa.PublicKey{N: big.NewInt(0)}
		pubKey.N.SetString(string(line), 10)
		line, _, err = buf.ReadLine(); checkErr(err)
		pubKey.E, err = strconv.Atoi(string(line)); checkErr(err)
		return &pubKey
		} else{
		fmt.Println("Error: unknown command", string(line))
		os.Exit(1)
	}
	return nil
}

func getRemoteConn(c *net.TCPConn)*remoteConn{
	return &remoteConn{c: c, pubK: waitPubKey(bufio.NewReader(c))}
}

func(rConn *remoteConn) sendCommand(comm string){
	eComm, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, rConn.pubK, []byte(comm), nil)
	checkErr(err)
	rConn.c.Write(eComm)
}

func listen(){
	l, err := net.ListenTCP(tcpProtocol, listenAddr)
	checkErr(err)
	fmt.Println("Listen port:", l.Addr().(*net.TCPAddr).Port)
	for {
		c, err := l.AcceptTCP()
		checkErr(err)
		fmt.Println("Connect from:", c.RemoteAddr())
		rConn := getRemoteConn(c)

		mes := ""
		for {
			fmt.Scanf("%s", &mes)
			rConn.sendCommand(mes)
			if mes == "&" {					//break command should be changed here as well as in client in main func
				break
			}
		}
	}
}

func main(){
	listen()
}