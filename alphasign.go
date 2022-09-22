package main

/*
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$1BL000E0000AAU0040FFFF1AU0064FEFE2AU0064FEFE3AU0064FEFE4AU0064FEFE5AU0064FEFE\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E.TUA\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02G1\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02AA\x1b bTesting 1-2-3\x04'
*/

import (
  "log"
  "encoding/hex"
)

const (
  disp_pos_middle = byte(0x20)
  disp_pos_top = byte(0x22)
  disp_pos_bottom = byte(0x26)
  disp_pos_fill = byte(0x30)
  disp_pos_left = byte(0x31)
  disp_pos_right = byte(0x32)

  disp_mode_rotate = 'a'
  disp_mode_hold = byte('b')
  disp_mode_flash = 'c'
  disp_mode_roll_up = 'e'
  disp_mode_roll_down = 'f'
  disp_mode_roll_left = 'g'
  disp_mode_roll_right = 'h'
  disp_mode_wipe_up = 'i'
  disp_mode_wipe_down = 'j'
  disp_mode_wipe_left = 'k'
  disp_mode_wipe_right = 'l'
  disp_mode_scroll = 'm'
  disp_mode_auto = 'o'
  disp_mode_roll_in = 'p'
  disp_mode_roll_out = 'q'
  disp_mode_wipe_in = 'r'
  disp_mode_wipe_out = 's'

  func_clear_set_mem = byte('$')
  func_set_run_seq = byte('.')

  file_type_text = byte('A')
  file_type_string = byte('B')
  file_type_pict = byte('C')
 
  file_prot_unlocked = byte('U')
  file_prot_locked = byte('L')

)

type memCfg struct {
  fileLabel byte
  fileType byte
  fileProt byte
  size int
  q int
}

var  eot = []byte{0x04}

func WriteSpecialFunction(function byte) []byte {
  var p = header()

  h  := []byte {'E', function}
  return append(p,h...)
}

func ClearMemory() []byte {
  p := WriteSpecialFunction(func_clear_set_mem)
  return append(p,eot...)
}
func SetMemory(mem []memCfg) []byte {
  p := WriteSpecialFunction(func_clear_set_mem)
  for _,m := range mem {
    log.Printf("%v",m)
  }
  return append(p,eot...)
}

func SetRunSequence() []byte {
  p:= WriteSpecialFunction(func_set_run_seq)
  return p
}

func ClearText() []byte{
  var p = header()
  // This appears to clear file label "1"?
  h  := []byte {'G',  '1'}
  p = append(p,h...)
  p = append(p,eot...)
  return (p)
}


func WriteText(text string ) []byte{
  var p = header()

  dispPos := disp_pos_middle
  mode :=  disp_mode_hold
  h  := []byte {'A', 'A', 0x1b, dispPos, mode}
  ba := []byte(text)
  for i,b := range text {
    if ((b < 32) || (b >= 127)) {
      ba[i] = byte('?')
    } else {
      ba[i] = byte(b)
    }
  }
  p = append(p,h...)
  p = append(p,ba...)
  p = append(p,eot...)
  return p
}


func header() []byte {
  var h  = []byte {0x00, 0x00, 0x00, 0x00, 0x01, 'Z', 0x00, 0x02}
  return h
}

func alphasign() {

  var packet = ClearText()
  log.Printf("ClearText:\n%s",hex.Dump(packet))
  packet = WriteText("Test")
  log.Printf("WriteText:\n%s",hex.Dump(packet))
}
