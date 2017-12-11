package main

import (
	"bufio"
	"fmt"
  "io"
  "strconv"
  "strings"
	"os"
)

var labels map[string]int

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func Whitespace(r rune) bool {
    return r == ' ' || r == '\t'
}

func lpad(s string,pad string, plength int)string{
    for i:=len(s);i<plength;i++{
        s=pad+s
    }
    return s
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
  fmt.Print(line)
  if strings.HasPrefix(line, ";") { //It's a comment
    return -88888888
  } else if len(line) == 1 {
    return -88888888
  } else {
    tokens := strings.SplitN(line, ":", 2)
    inst := tokens[0]
    if len(tokens) > 1 { //There is a label
      labels[tokens[0]] = addr;
      inst = tokens[1]
    }
    inst = strings.TrimFunc(inst, Whitespace)
    args := strings.FieldsFunc(inst, Whitespace)
    opcode := strings.TrimFunc(args[0], Whitespace)
    if opcode == "DS" {
      if len(args) > 1 {
        res, err := strconv.ParseInt(args[1], 10, 20)
        if err != nil {
          return 0
        }
        return int(res)
      } else {
        return 0
      }
    } else if len(args) > 2 && args[1] == "=" {
      res, err := strconv.ParseInt(args[2], 10, 20)
      check(err)
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

  f, err := os.Open(os.Args[1])
  check(err)

  r := bufio.NewReader(f)

  for addr := 0; true; addr++ {
    line, err := r.ReadString('\n')
    if err == io.EOF {
      break;
    }

    res := parseLine(line, addr)
    if res == -88888888 {
      addr--
    } else {
      res = res & 33554431
      fmt.Printf("%024b\n0x%06X\n", res, res)
    }
  }
  fmt.Println(labels)
}
