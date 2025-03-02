package main

/*
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$1BL000E0000AAU0040FFFF1AU0064FEFE2AU0064FEFE3AU0064FEFE4AU0064FEFE5AU0064FEFE\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E.TUA\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02G1\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02AA\x1b bTesting 1-2-3\x04'


root@speakbot:/home/bkg# stty -a -F /dev/ttyUSB0
speed 4800 baud; rows 0; columns 0; line = 0;
intr = ^C; quit = ^\; erase = ^?; kill = ^U; eof = ^D; eol = <undef>; eol2 = <undef>; swtch = <undef>; start = ^Q; stop = ^S; susp = ^Z; rprnt = ^R; werase = ^W; lnext = ^V; discard = ^O;
min = 0; time = 0;
parenb -parodd -cmspar cs7 hupcl cstopb cread clocal -crtscts
-ignbrk -brkint -ignpar -parmrk -inpck -istrip -inlcr -igncr -icrnl -ixon -ixoff -iuclc -ixany -imaxbel -iutf8
-opost -olcuc -ocrnl -onlcr -onocr -onlret -ofill -ofdel nl0 cr0 tab0 bs0 vt0 ff0
-isig -icanon -iexten -echo -echoe -echok -echonl -noflsh -xcase -tostop -echoprt -echoctl -echoke -flusho -extproc



   self._conn = serial.Serial(port=self.device,
                               baudrate=4800,
                               parity=serial.PARITY_EVEN,
                               stopbits=serial.STOPBITS_TWO,
                               bytesize=serial.SEVENBITS,
                               timeout=1,
                               xonxoff=0,
                               rtscts=0)


4800 / 7 / 2 / Even
*/

import (
  "log"
  //"encoding/hex"
  "os"
  "bytes"
  "os/exec"
  "fmt"
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

  run_seq_time = byte('T')
  run_seq_order = byte('S')
  run_seq_delete = byte('D')

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
    b:= []byte{m.fileLabel,m.fileType,m.fileProt}
    p =  append(p,b...)
    x := fmt.Sprintf("%04X",m.size)
    p =  append(p,[]byte(x)...)
    x = fmt.Sprintf("%04X",m.q)
    p =  append(p,[]byte(x)...)
  }
  return append(p,eot...)
}

func SetRunSequence() []byte {
  p:= WriteSpecialFunction(func_set_run_seq)

  b:= []byte{run_seq_time, file_prot_unlocked,'A'}
  p =  append(p,b...)
  return append(p,eot...)
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
  var h  = []byte {0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 'Z', byte('0'), byte('0'), 0x02}
  return h
}

func alphasign_init(port string) {
  args := []string{"4800", "cs7", "parenb", "raw", "-parodd", "-crtscts", "-echo", ("-F" + port)}
  cmd := exec.Command("/bin/stty", args...)
  var errout bytes.Buffer
  cmd.Stderr = &errout
  err := cmd.Run()
  
  if err != nil {
    log.Fatal("Error alphasign init Error: ", errout.String())
    return
  } 
  cmd = exec.Command("/bin/echo", args...)
  err = cmd.Run()
}

func alphasign(text string,port string) {

  alphasign_init(port)
  //fd,err := os.Open("/dev/ttyUSB0")
  fd, err := os.OpenFile("/dev/ttyUSB0", os.O_APPEND|os.O_WRONLY, 0644)
  defer fd.Close()
  if (err != nil) { log.Fatal("Error opening port",err) }
  var packet = ClearMemory()
  //log.Printf("ClearMemory:\n%s",hex.Dump(packet))
  fd.Write(packet)

  mem := []memCfg{
    {'1',file_type_string,file_prot_locked,0x000E,0x0000},
    {'A',file_type_text,file_prot_unlocked,0x0040,0xFFFF},
    {'1',file_type_text,file_prot_unlocked,0x0064,0xFEFE},
    {'2',file_type_text,file_prot_unlocked,0x0064,0xFEFE},
    {'3',file_type_text,file_prot_unlocked,0x0064,0xFEFE},
    {'4',file_type_text,file_prot_unlocked,0x0064,0xFEFE},
    {'5',file_type_text,file_prot_unlocked,0x0064,0xFEFE}}
  packet = SetMemory(mem)
  //log.Printf("SetMemory:\n%s",hex.Dump(packet))
  fd.Write(packet)


  packet = SetRunSequence()
  //log.Printf("RunSequence:\n%s",hex.Dump(packet))
  fd.Write(packet)

  packet = ClearText()
  //log.Printf("ClearText:\n%s",hex.Dump(packet))
  fd.Write(packet)

  packet = WriteText(text)
  //log.Printf("WriteText:\n%s",hex.Dump(packet))
  fd.Write(packet)
  
}
