package gotorrent

// A metainfo file (.torrent) gives info about a torrent file.
// See https://wiki.theory.org/BitTorrentSpecification#Metainfo_File_Structure for details.
type MetaInfo struct {
	Announce      string
	Announce_List [][]string
	Comment       string
	CreatedBy     string
	CreationDate  int
	Encoding      string
	InfoHash      string
	Info          struct {
		Name        string
		PieceLength int
		Pieces      string
		Length      int
		Files       []struct {
			Length int
			Path   []string
		}
	}
}
