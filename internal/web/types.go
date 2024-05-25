package web

type SnapInfo struct {
	FolderName string
	FileName   string
	FilePath   string
}

type TimelapseSnap struct {
	FileName string
	FilePath string
}
type Timelapse struct {
	FolderName string
	FolderPath string
	Snaps      []TimelapseSnap
}
type Timelapses []Timelapse

type TLInfo struct {
	FolderName        string
	NumberOfSnaps     uint
	HasTimelapseVideo bool
	FolderPath        string
}

type Hi struct {
	ThumbnailPath string
	ix            int
	ImgPath       string
}
