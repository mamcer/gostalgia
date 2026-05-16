package domain

import (
	"time"

	"gorm.io/datatypes"
)

type NFile struct {
	ID           int64          `json:"id" gorm:"primaryKey;column:id"`
	Name         string         `json:"name" gorm:"column:name"`
	Extension    string         `json:"extension" gorm:"column:extension"`
	Path         string         `json:"path" gorm:"column:path"`
	DateModified time.Time      `json:"date_modified" gorm:"column:date_modified"`
	CapturedAt   *time.Time     `json:"captured_at" gorm:"column:captured_at"`
	Size         int64          `json:"size" gorm:"column:size"`
	Importance   int            `json:"importance" gorm:"column:importance"`
	Hash         string         `json:"hash" gorm:"column:hash"`
	Metadata     datatypes.JSON `json:"metadata" gorm:"column:metadata"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	Tags         []*NTag        `json:"tags" gorm:"many2many:ntag_nfile;foreignKey:ID;joinForeignKey:nfile_id;References:ID;joinReferences:ntag_id"`
}

func (NFile) TableName() string {
	return "nfile"
}

type NDirectory struct {
	ID                int64        `json:"id" gorm:"primaryKey;column:id"`
	Name              string       `json:"name" gorm:"column:name"`
	DateModified      time.Time    `json:"date_modified" gorm:"column:date_modified"`
	ParentDirectoryID int64        `json:"parent_directory_id" gorm:"column:parent_directory_id"`
	ParentDirectory   *NDirectory  `json:"parent_directory" gorm:"foreignKey:ParentDirectoryID"`
	FullPath          string       `json:"full_path" gorm:"column:full_path"`
	Size              int64        `json:"size" gorm:"column:size"`
	FileCount         int64        `json:"file_count" gorm:"column:file_count"`
	DirectoryCount    int64        `json:"directory_count" gorm:"column:directory_count"`
	IsSource          bool         `json:"is_source" gorm:"column:is_source"`
	FileNodes         []*NFileNode `json:"file_nodes" gorm:"foreignKey:NDirectoryID"`
	Tags              []*NTag      `json:"tags" gorm:"many2many:ntag_ndirectory;foreignKey:ID;joinForeignKey:ndirectory_id;References:ID;joinReferences:ntag_id"`
}

func (NDirectory) TableName() string {
	return "ndirectory"
}

type NFileNode struct {
	ID           int64       `json:"id" gorm:"primaryKey;column:id"`
	Name         string      `json:"name" gorm:"column:name"`
	NFileID      int64       `json:"nfile_id" gorm:"column:nfile_id"`
	File         *NFile      `json:"file" gorm:"foreignKey:NFileID"`
	NScanID      int64       `json:"nscan_id" gorm:"column:nscan_id"`
	Scan         *NScan      `json:"scan" gorm:"foreignKey:NScanID"`
	NDirectoryID int64       `json:"ndirectory_id" gorm:"column:ndirectory_id"`
	Directory    *NDirectory `json:"directory" gorm:"foreignKey:NDirectoryID"`
}

func (NFileNode) TableName() string {
	return "nfilenode"
}

type NScanStatus int

const (
	NScanStatusNotStarted NScanStatus = 0
	NScanStatusInProgress NScanStatus = 1
	NScanStatusCompleted  NScanStatus = 2
	NScanStatusFailed     NScanStatus = 3
)

type NScan struct {
	ID                        int64       `json:"id" gorm:"primaryKey;column:id"`
	DateCreated               time.Time   `json:"date_created" gorm:"column:date_created"`
	Duration                  int64       `json:"duration" gorm:"column:duration"`
	FileCount                 int64       `json:"file_count" gorm:"column:file_count"`
	DirectoryCount            int64       `json:"directory_count" gorm:"column:directory_count"`
	ExistingFileRepeatedCount int64       `json:"existing_file_repeated_count" gorm:"column:existing_file_repeated_count"`
	InternalFileRepeatedCount int64       `json:"internal_file_repeated_count" gorm:"column:internal_file_repeated_count"`
	Status                    NScanStatus `json:"status" gorm:"column:status"`
	RootDirectoryID           int64       `json:"root_directory_id" gorm:"column:root_directory_id"`
	RootDirectory             *NDirectory `json:"root_directory" gorm:"foreignKey:RootDirectoryID"`
}

func (NScan) TableName() string {
	return "nscan"
}

type NTag struct {
	ID          int64         `json:"id" gorm:"primaryKey;column:id"`
	Name        string        `json:"name" gorm:"column:name"`
	Files       []*NFile      `json:"files" gorm:"many2many:ntag_nfile;foreignKey:ID;joinForeignKey:ntag_id;References:ID;joinReferences:nfile_id"`
	Directories []*NDirectory `json:"directories" gorm:"many2many:ntag_ndirectory;foreignKey:ID;joinForeignKey:ntag_id;References:ID;joinReferences:ndirectory_id"`
}

func (NTag) TableName() string {
	return "ntag"
}
