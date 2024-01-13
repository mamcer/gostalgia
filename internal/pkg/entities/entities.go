package entities

import "time"

const (
	Done int64 = iota + 1
	InProgress
	Error
)

type Nscan struct {
	ID                int64     // scan id
	Name              string    // scan name
	DateCreated       time.Time // scan creation date
	Duration          int64     // scan duration (in milliseconds)
	FileCount         int64     // file scan count
	DirectoryCount    int64     // directory scan count
	FileRepeatedCount int64     // file repeated scan count
	Status            int64     // scan status = done, inprogress, error
	RootDirectoryId   int64     // scan directory id
	RetryCount        int64     // scan retry count
}

type Ndirectory struct {
	ID             int64     // directory id
	Name           string    // directory name
	Path           string    // directory path
	DateModified   time.Time // date modified
	Size           int64     // directory size (in bytes)
	FileCount      int64     // directory file count
	Fpath          string    // current file path
	DirectoryCount int64     // directory directory count
	ParentID       int64     // parent directory id
}

type Nfile struct {
	ID           int64     // file id
	Name         string    // file name
	Extension    string    //file extension
	Path         string    // file path
	DateModified time.Time // file date modified
	Size         int64     // file size (in bytes)
	Hash         string    // file hash
}

type Nfilendirectory struct {
	ID           int64  // file scan id
	NfileID      int64  // file id
	NdirectoryID int64  // directory id
	NscanID      int64  // scan id
	Name         string //file name
}

type Nerror struct {
	ID          int64  // error id
	Description string // error description
	NscanID     int64  // scan id
}
