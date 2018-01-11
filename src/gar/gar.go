package gar

import (
	"errors"
	"fmt"
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
	ar.headList = make([]*fileHeader, 0, 16)

	//br := bufio.NewReader(f)
	buf := make([]byte, 60)
	n, err := f.Read(buf[:8])
	//	line, _ := br.ReadString('\n')

	if n == 8 && strings.Compare("!<arch>\n", string(buf[:8])) != 0 {
		return nil, errors.New("not ar file!")
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
		ar.headList = append(ar.headList, head)

		// read data
		//pad := head.size % 2
		//f.Seek(int64(head.size+pad), os.SEEK_CUR)
		head.data = make([]byte, head.size)
		f.Read(head.data)
		if head.size%2 == 1 {
			f.Seek(1, os.SEEK_CUR)
		}
	}
	return ar, nil
}
func (arf *ArFile) List() {
	for _, v := range arf.headList {
		fmt.Println(v)
	}
}
func (arf *ArFile) Append(file string) {

}
func (arf *ArFile) Close() {

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
