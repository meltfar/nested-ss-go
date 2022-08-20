package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var config struct {
	Verbose bool
	TCPCork bool
}

func main() {
	var flags struct {
		ServerAddr  string
		RealBackend string
		Cipher      string
		Key         string
		Password    string
		Keygen      int
		Socks       string
	}

	flag.BoolVar(&config.Verbose, "verbose", false, "verbose mode")
	flag.StringVar(&flags.Cipher, "cipher", "chacha20-ietf-poly1305", "available ciphers: "+"aes-256-gcm chacha20-ietf-poly1305")
	flag.StringVar(&flags.Key, "key", "", "base64url-encoded key (derive from password if empty)")
	flag.StringVar(&flags.RealBackend, "real", "", "real backend ss server address (with port)")
	flag.IntVar(&flags.Keygen, "keygen", 0, "generate a base64url-encoded random key of given length in byte")
	flag.StringVar(&flags.Password, "password", "Jfyytd000.", "password")
	flag.StringVar(&flags.ServerAddr, "n", "127.0.0.1:7890", "next socks5 address to proxy on")
	flag.StringVar(&flags.Socks, "l", "127.0.0.1:8688", "(client-only) SOCKS listen address")
	flag.BoolVar(&config.TCPCork, "tcpcork", false, "coalesce writing first few packets")
	flag.Parse()

	if flags.Keygen > 0 {
		key := make([]byte, flags.Keygen)
		_, _ = io.ReadFull(rand.Reader, key)
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		return
	}

	if flags.ServerAddr == "" {
		flag.Usage()
		return
	}

	var key []byte
	if flags.Key != "" {
		k, err := base64.URLEncoding.DecodeString(flags.Key)
		if err != nil {
			log.Fatal(err)
		}
		key = k
	}

	ssServerAddr := flags.ServerAddr
	cipher := flags.Cipher
	password := flags.Password
	var err error

	if strings.HasPrefix(ssServerAddr, "ss://") {
		ssServerAddr, cipher, password, err = parseURL(ssServerAddr)
		if err != nil {
			log.Fatal(err)
		}
	}

	realBack := flags.RealBackend

	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}

	if flags.Socks != "" {
		socks.UDPEnabled = false
		go socksLocal(flags.Socks, ssServerAddr, realBack, ciph.StreamConn)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
