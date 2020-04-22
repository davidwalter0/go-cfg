// deepcopy gen not currently in use

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-flag"
	yaml "gopkg.in/yaml.v3"
)

// APP ...
// +k8s:deepcopy-gen=true
type APP struct {
	A struct {
		B struct {
			C struct {
				D struct {
					E int ``
				}
			}
		}
		IntegerValueOfX int
	}
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
	Outer
	UnsetUserArray []string
	UnsetIntArray  []int64
}

// +k8s:deepcopy-gen=true
type Inner struct {
	I   int                `default:"42"`
	F   float64            `default:"3.1415926"`
	Msi map[string]int     `default:"white:0,black:1,red:,green:2,blue:3"`
	Msf map[string]float64 `default:"e:2.71828,π:3.14159,ξ:1,ρ:.01,φ:1.2,β:3,α:.01"`
}

// +k8s:deepcopy-gen=true
type Outer struct {
	I int     `default:"42"`
	F float64 `default:"3.1415926"`
	Inner
	Msi map[string]int `default:"white:0,black:1,red:,green:2,blue:3"`
}

// APP1 ...
// +k8s:deepcopy-gen=true
type APP1 struct {
	APP
}

// APP2 ...
// +k8s:deepcopy-gen=true
type APP2 struct {
	A struct {
		B struct {
			C struct {
				D struct {
					E int ``
				}
			}
		}
		IntegerValueOfX int
	}
	Required     uint32        `short:"required" usage:"uint32 test" required:"1"`
	UserArray    []string      `default:"a,b,c,9,8,7"`
	IntArray     []int64       `default:"4,5,6"`
	Debug        bool          `name:"Debug" short:"d" default:"true" usage:"enable debug mode"`
	Port         int           `short:"p" default:"8888" usage:"primary ip port"`
	CaC          string        `usage:"cac for ..." default:"abc123"`
	CC           string        `usage:"cc users for ..."`
	User         string        `usage:"user for ..."`
	UserName     string        `name:"USER_NAME"`
	Users        []string      `name:"nameOverride"`
	Rate         float64       `default:"2.71828"`
	RateOfTravel float32       `short:"rt"  default:"3.1"`
	Timeout      time.Duration `short:"t1"  default:"24h1m3s"`
	Timeout2     time.Duration `short:"t2"  default:"24h1m3s" usage:"Timeout2 something"`
	Int8         int8          `short:"i8"  default:"-1"      usage:"int8   test"`
	Nint8        int8          `short:"n8"  default:"-2"     usage:"nint8  test"`
	Uint8        uint8         `short:"u8"  default:"3"      usage:"uint8  test"`
	Int16        int16         `short:"i16" default:"4"    usage:"int16  test"`
	Nint16       int16         `short:"n16" default:"-2"   usage:"nint16 test"`
	Uint16       uint16        `short:"u16" default:"5"    usage:"uint16 test"`
	Int32        int32         `short:"i32" default:"1048576"  usage:"int32  test"`
	Nint32       int32         `short:"n32" default:"-12"    usage:"nint32 test"`
	Uint32       uint32        `short:"u32" usage:"uint32 test" required:"1"`
	// Uint32         uint32             `short:"u32" default:"255"      usage:"uint32 test" required:""`
	ColorCodes map[string]int     `short:"color"  default:"white:0xfff,black:000,red:f00,green:0f0,blue:00f"`
	Map        map[string]float64 `short:"m"   default:"π:3.14159,ξ:1,ρ:.01,φ:1.2,β:3,α:.01,δ:3,ε:.001,φ:.1,ψ:.9,ω:2.1"`
	Outer
	UnsetUserArray []string
	UnsetIntArray  []int64
}

// +k8s:deepcopy-gen=true
type Flags struct {
	Help bool
	Load bool
	Save bool
	File string
}

func main() {

if false {
	flag.NonFlagOSArgs(1)
}
	p := &APP{}
	p1 := &APP1{}
	p2 := &APP2{}
	flg := &Flags{}
	cfg.NestWrap("w", p, p1, p2)
	cfg.Flags(flg)
	fmt.Println("os.Args", os.Args)
	// cfg.Usage()
	// fmt.Println(os.Args)
	// for _, f := range os.Args {
	// 	if strings.Index(f, "--help") == 0 {
	// 		cfg.Usage()
	// 	}
	// }
	// cfg.Usage()
	fmt.Printf("%+v\n", flg)
	if flg.Help {
		cfg.Usage()
		os.Exit(1)
	}
	if flg.Save {
		cfg.Store.Save("config-2.yaml")
	}
	if flg.Load {
		cfg.Store.Load("config.yaml")
    // fmt.Println("Dump Store",Dump(cfg.Store))
		cfg.Store.Save(flg.File)
	}

	fmt.Println(Dump(flg))
	d, err := time.ParseDuration("720h1m3s")
	fmt.Printf("d %s %+v %v\n", d, d, err)
	fmt.Println("Dump d", d)
  fmt.Println("Dump p",Dump(p))
  fmt.Println("Dump p1",Dump(p1))
  fmt.Println("Dump Store",Dump(cfg.Store))
	// for _, f := range os.Args {
	// 	if strings.Index(f, "--help") == 0 {
	// 		flag.Parse()
	// 		flag.Usage()
	// 		cfg.Usage()
	// 		os.Exit(1)
	// 	}
	// }

}
func Dump(o interface{}) string {
	byte, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(byte))

	byte, err = yaml.Marshal(o)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(byte))
	return string(byte)
}
