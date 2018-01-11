package gar

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func init() {
	fmt.Println("init")
}

type ArFile struct {
	file     string
	headList []*fileHeader
	modified bool
}

type fileHeader struct {
	id        string //16
	timestamp string
	ownerId   string
	groupId   string
	mode      string
	size      int
	offset    int64
	data      []byte
}

func (h *fileHeader) String() string {
	return "\t" + h.id + " " + strconv.Itoa(h.size)
}
func Open(file string) (ar *ArFile, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ar = &ArFile{}
	ar.file = file

	return ar, nil
}
func (arf *ArFile) List() []string {
	f, err := os.Open(arf.file)
	if err != nil {
		return nil
	}
	defer f.Close()

	headList := make([]*fileHeader, 0, 16)

	//br := bufio.NewReader(f)
	buf := make([]byte, 60)
	n, err := f.Read(buf[:8])
	//	line, _ := br.ReadString('\n')

	if n == 8 && strings.Compare("!<arch>\n", string(buf[:8])) != 0 {
		return nil
	}
	for {
		//line, _ = br.ReadString('\n')
		n, err = f.Read(buf)
		if err != nil {
			break
		}

		head := parseHead(string(buf))
		offset, _ := f.Seek(0, os.SEEK_CUR)
		head.offset = offset
		headList = append(headList, head)

		// read data
		pad := head.size % 2
		f.Seek(int64(head.size+pad), os.SEEK_CUR)
		/*head.data = make([]byte, head.size)
		f.Read(head.data)
		if head.size%2 == 1 {
			f.Seek(1, os.SEEK_CUR)
		}*/
	}

	for _, v := range headList {
		fmt.Println(v)
	}
	return nil
}
func (arf *ArFile) Append(file string) error {
	f, err := os.Open(arf.file)
	if err != nil {
		return err
	}
	defer f.Close()
	//add header
	f.Seek(0, os.SEEK_END)
	var head [60]byte
	w := bytes.NewBuffer(head[:])
	id := f.Name()
	if len(id) > 16 {
		id = id[:16]
	}
	w.WriteString(id)
	// write space
	f.Write(head[:])
	data, err := ioutil.ReadFile(file)
	// add data
	f.Write(data)
	pad := []byte{'\n'}
	if len(data)%2 == 1 {
		f.Write(pad)
	}
	return nil
}
func (arf *ArFile) Close() {
	if arf.modified {
		//write data back
	}
}

func parseHead(line string) *fileHeader {
	head := &fileHeader{}
	head.id = strings.TrimRight(line[0:16], "/ ")
	head.timestamp = strings.TrimRight(line[16:28], " ")
	head.ownerId = strings.TrimRight(line[28:34], " ")
	head.groupId = strings.TrimRight(line[34:40], " ")
	head.mode = strings.TrimRight(line[40:48], " ")
	head.size, _ = strconv.Atoi(strings.TrimRight(line[48:58], " "))

	return head
}
