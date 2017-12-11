package main

import (
	"bufio"
  "flag"
	"fmt"
  "io"
  "strconv"
  "strings"
	"os"
)

var labels map[string]int

func check(e error, errMsg string) {
  if e != nil {
    panic("P:" + errMsg + ":" + e.Error())
  }
}

func Whitespace(r rune) bool {
    return r == ' ' || r == '\t'
}

func getOpcode(inst string) (int) {
  if len(inst) >= 2 {
    if strings.HasPrefix(inst, "LD") {
      switch inst[2] {
      case 'C': //LDC
        return 0
      case 'V': //LDV
        return 1
      case 'I': //LDIV
        return 0xA
      }
    } else if strings.HasPrefix(inst, "ST") {
      switch inst[2] {
      case 'V': //STV
        return 2
      case 'I': //STIV
        return 0xB
      }
    } else if strings.HasPrefix(inst, "A") {
      switch inst[1] {
      case 'D': //ADD
        return 3
      case 'N': //AND
        return 5
      }
    } else if strings.HasPrefix(inst, "OR") {
      return 5
    } else if strings.HasPrefix(inst, "XOR") {
      return 6
    } else if strings.HasPrefix(inst, "EQL") {
      return 7
    } else if strings.HasPrefix(inst, "JI") {
      return 13
    } else if strings.HasPrefix(inst, "JM") {
      switch inst[2] {
      case 'P': //JMP
        return 8
      case 'N': //JMN
        return 9
      case 'S': //JMS
        return 12
      }
    } else if strings.HasPrefix(inst, "HALT") {
      return 0xF0
    } else if strings.HasPrefix(inst, "NOT") {
      return 0xF1
    } else if strings.HasPrefix(inst, "RAR") {
      return 0xF2
    } else {
      panic("Invalid operation")
    }
  } else {
    panic("Invalid operation")
  }
  return -88888888
}

func parseLine(line string, addr int) (int) {
  line = strings.TrimFunc(line, Whitespace)
  if strings.HasPrefix(line, ";") { //It's a comment
    return -88888888
  } else if len(line) == 1 {
    return -88888888
  } else {
    tokens := strings.SplitN(line, ":", 2)
    inst := tokens[0]
    if len(tokens) > 1 { //There is a label
      fmt.Printf("P:INFO:Setting label %s at 0x%06X\n", tokens[0], addr)
      labels[tokens[0]] = addr;
      inst = tokens[1]
    }
    inst = strings.TrimFunc(inst, Whitespace)
    args := strings.FieldsFunc(inst, Whitespace)
    opcode := strings.TrimFunc(args[0], Whitespace)
    if strings.HasPrefix(opcode, "DS") {
      if len(args) > 1 {
        res, err := strconv.ParseInt(strings.Trim(args[1], "\n"), 10, 20)
        if err != nil {
          fmt.Printf("P:WARN:Parsing DS parameter failed: %s\n", err.Error())
          fmt.Printf("P:INFO:DS without parameter, initializing at 0\n")
          return 0
        }
        fmt.Printf("P:INFO:Setting 0x%06X to %d\n", addr, res)
        return int(res)
      } else {
        return 0
      }
    } else if len(args) > 2 && args[1] == "=" {
      res, err := strconv.ParseInt(args[2], 10, 20)
      check(err, fmt.Sprintf("ERR:Error parsing argument for global at 0x%06X", addr))
      labels[args[0]] = int(res)
      return -88888888
    } else {
      code := getOpcode(opcode)
      if code < 16 {
        res, err := strconv.ParseInt(args[1], 10, 20)
        if err != nil {
          if val, ok := labels[args[1]]; ok {
            res = int64(val)
          } else {
            res = 0
          }
        }

        return code * 1048576 + int(res)
      } else {
        return code * 65536
      }
    }
  }
}

func main() {
  labels = make(map[string]int)

  listing := ""

  magicNum := []byte{74, 66}

  var binListing = flag.Bool("bindump", false, "Create listing as binary instead of hexadecimal")
  var outFile = flag.String("out", "out.mib", "Binary output file name")

  flag.Parse()

  f, err := os.Open(flag.Arg(0))
  check(err, "ERR:Error opening file")

  r := bufio.NewReader(f)

  fout, err := os.Create(*outFile)
  check(err, "ERR:Error creating output file")

  _, err = fout.Write(magicNum)
  check(err, "ERR:Error writing to output file")

  fout.Sync()

  for addr := 0; true; addr++ {
    line, err := r.ReadString('\n')
    if err == io.EOF {
      break;
    } else if err != nil {
      check(err, "ERR:Error reading source code file")
    }

    res := parseLine(line, addr)
    if res == -88888888 {
      addr--
    } else {
      res = res & 16777215
      if *binListing {
        listing += fmt.Sprintf("%024b %024b\n", addr, res)
      } else {
        listing += fmt.Sprintf("0x%06X 0x%06X\n", addr, res)
      }

      fout.WriteString(fmt.Sprintf("%024b", res))
      fout.Sync()
    }
  }
  fmt.Println("Listing:\n\n" + listing)
}
