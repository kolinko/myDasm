package main

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
)

func readCSVFile(input string) (data [][]string, err error) {
	sf, err := os.Open(input)
	if err != nil {
		fmt.Println(err, "[", input, "]")
		log.Fatal(err, input)
	}
	data, err = csv.NewReader(sf).ReadAll()
	return
}

type replyRec struct {
	ID             uint64
	CreatedAt      string `json:"created_at"`
	TextSignature  string `json:"text_signature"`
	HexSignature   string `json:"hex_signature"`
	BytesSignature string `json:"bytes_signature"`
}

type reply struct {
	Count    int
	Next     string
	Previous string
	Results  []replyRec
}

type opData struct {
	Gas         int64
	Description string
}

func gas(op vm.OpCode) string {
	if !showGas {
		return ""
	}
	res := opCodeMap[byte(op)].Gas
	if res == -1 {
		return "(*)"
	}
	return fmt.Sprintf("(%d)", res)
}

func opcodeDescription(op vm.OpCode) string {
	if !showOpDesc {
		return ""
	}
	return opCodeMap[byte(op)].Description

}

func printDisassembled(code string) error {
	script, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	it := asm.NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			if len(it.Arg()) == 4 && it.Op() == vm.PUSH4 {
				call := fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/?hex_signature=0x%08x", it.Arg())
				comment := ""
				repl, err := http.Get(call)
				if err == nil {
					var b4 reply
					err = json.NewDecoder(repl.Body).Decode(&b4)
					if err == nil {
						for _, b4rec := range b4.Results {
							comment += " " + b4rec.TextSignature
						}
					}
				}
				if len(comment) > 0 {
					comment = " // " + comment
				}
				fmt.Printf("%06X: %s %v 0x%x %s\n", it.PC(), gas(it.Op()), it.Op(), it.Arg(), comment)

			} else if it.Op() == vm.PUSH32 {
				ascii := ""
				pos := bytes.IndexByte(it.Arg(), 0)
				ok := true
				for i := 0; i < pos; i++ {
					ok = ok && it.Arg()[i] < 0x80
				}
				if ok {
					if pos < 0 {
						ascii = " // " + string(it.Arg())
					} else {
						ascii = " // " + string(it.Arg()[:pos])
					}
				}
				fmt.Printf("%06X: %s %v 0x%x %s\n", it.PC(), gas(it.Op()), it.Op(), it.Arg(), ascii)

			} else {
				fmt.Printf("%06X: %s %v 0x%x %s\n", it.PC(), gas(it.Op()), it.Op(), it.Arg(), opcodeDescription(it.Op()))
			}
		} else {
			fmt.Printf("%06X: %s %v %s\n", it.PC(), gas(it.Op()), it.Op(), opcodeDescription(it.Op()))
		}
	}
	return it.Error()
}

var input string
var opCodeMap map[byte]opData
var showGas bool
var showOpDesc bool

func main() {
	opCodeMap = make(map[byte]opData)
	flag.StringVar(&input, "input", "", " path to input file (Hex)")
	flag.BoolVar(&showGas, "gas", false, "diplay gas usage per opcode")
	flag.BoolVar(&showOpDesc, "opcode", false, "display opcode function as comment")
	flag.Parse()
	if len(input) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	gasData, err := readCSVFile("gas.csv")
	if err == nil {
		for _, line := range gasData {
			pos := common.FromHex(line[0])
			if len(pos) == 0 {
				log.Fatal("Zero length position ", line)
			}
			if len(pos) > 1 {
				log.Fatalf("Error - opcode > 256 (%d)", len(pos))
			}
			od := opData{}
			if len(line[2]) > 0 {
				od.Description = "// " + line[2]
			}
			if num, err := strconv.ParseInt(line[4], 10, 16); err == nil {
				od.Gas = num
			} else {
				od.Gas = -1
			}
			//fmt.Println(pos[0], od.Description, od.Gas)
			opCodeMap[pos[0]] = od
		}
	}

	in, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	code := strings.TrimSpace(string(in[:]))
	//fmt.Printf("%v\n", code)
	err = printDisassembled(code)
	if err != nil {
		log.Fatal(err)
	}
}
