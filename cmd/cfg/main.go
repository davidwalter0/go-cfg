package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidwalter0/go-cfg"
)

type MyAPP struct {
	Required     uint32        `short:"required" usage:"uint32 test" required:"1"`
	UserArray    []string      `default:"x,y,z,0,1"`
	IntArray     []int64       `default:"0,1,2,3,4"`
	Debug        bool          `name:"Debug" short:"d" default:"false" usage:"enable debug mode"`
	Port         int           `short:"p" default:"8080" usage:"primary ip port"`
	CaC          string        `usage:"cc users for ..." default:"abc123"`
	CC           string        `usage:"cc users for ..."`
	User         string        `usage:"user for ..."`
	UserName     string        `name:"USER_NAME"`
	Users        []string      `name:"nameOverride"`
	Rate         float64       `default:"2.71828"`
	RateOfTravel float32       `short:"rt"  default:"3.14"`
	Timeout      time.Duration `short:"t1"  default:"720h1m3s"`
	Timeout2     time.Duration `short:"t2"  default:"720h1m3s" usage:"Timeout2 something"`
	Int8         int8          `short:"i8"  default:"127"      usage:"int8   test"`
	Nint8        int8          `short:"n8"  default:"-128"     usage:"nint8  test"`
	Uint8        uint8         `short:"u8"  default:"255"      usage:"uint8  test"`
	Int16        int16         `short:"i16" default:"32767"    usage:"int16  test"`
	Nint16       int16         `short:"n16" default:"-32768"   usage:"nint16 test"`
	Uint16       uint16        `short:"u16" default:"65535"    usage:"uint16 test"`
	Int32        int32         `short:"i32" default:"1048576"  usage:"int32  test"`
	Nint32       int32         `short:"n32" default:"-1232"    usage:"nint32 test"`
	Uint32       uint32        `short:"u32" usage:"uint32 test" required:"1"`
	// Uint32         uint32             `short:"u32" default:"255"      usage:"uint32 test" required:""`
	ColorCodes map[string]int     `short:"color"  default:"white:0xfff,black:000,red:f00,green:0f0,blue:00f"`
	Map        map[string]float64 `short:"m"   default:"π:3.14159,ξ:1,ρ:.01,φ:1.2,β:3,α:.01,δ:3,ε:.001,φ:.1,ψ:.9,ω:2.1"`
	Outer      struct {
		I     int            `default:"42"`
		F     float64        `default:"3.1415926"`
		Msi   map[string]int `default:"white:0,black:1,red:,green:2,blue:3"`
		Inner struct {
			I   int                `default:"42"`
			F   float64            `default:"3.1415926"`
			Msi map[string]int     `default:"white:0,black:1,red:,green:2,blue:3"`
			Msf map[string]float64 `default:"e:2.71828,π:3.14159,ξ:1,ρ:.01,φ:1.2,β:3,α:.01"`
		}
	}
	A struct {
		B struct {
			C struct {
				D struct {
					E int ``
				}
			}
		}
	}
	UnsetUserArray []string
	UnsetIntArray  []int64
}

func main() {
	p := &MyAPP{}
	// c := cfg.NewConfigurer(p)
	cfg.Eval(p)

	byte, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(byte))

	for _, f := range os.Args {
		if strings.Index(f, "--help") == 0 {
			flag.Parse()
			flag.Usage()
			cfg.Usage()
			os.Exit(1)
		}
	}

}
