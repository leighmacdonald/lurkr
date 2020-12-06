package torrent

import (
	"github.com/anacrolix/torrent/bencode"
	"io"
)

type File struct {
	Info struct {
		Name        string `bencode:"name"`
		Length      int64  `bencode:"length"`
		MD5Sum      string `bencode:"md5sum,omitempty"`
		PieceLength int64  `bencode:"piece length"`
		Pieces      string `bencode:"pieces"`
		Private     bool   `bencode:"private,omitempty"`
	} `bencode:"info"`

	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list,omitempty"`
	CreationDate int64       `bencode:"creation date,omitempty"`
	Comment      string      `bencode:"comment,omitempty"`
	CreatedBy    string      `bencode:"created by,omitempty"`
	URLList      interface{} `bencode:"url-list,omitempty"`
}

func (f File) Size() uint64 {
	panic("Calculate torrent size")
	return 0
}

func Decode(reader io.Reader, file *File) error {
	return bencode.NewDecoder(reader).Decode(file)
}

func Encode(writer io.Writer, file *File) error {
	return bencode.NewEncoder(writer).Encode(file)
}
