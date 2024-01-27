# Nostalgia

2024-01-14


despues de unas 7hs corriendo. Todo el contenido del disco

47.32% - [new] /media/darkforce/documents/backups/yani-facultad/Facultad/Orange/www.ailab.si/orange/acknowledgements.htm
47.32% - [new] /media/darkforce/documents/backups/yani-facultad/Facultad/Orange/www.ailab.si/orange/buildC45.py
47.32% - [new] /media/darkforce/documents/backups/yani-facultad/Facultad/Orange/www.ailab.si/orange/customInstall.htm
47.32% - [new] /media/darkforce/documents/backups/yani-facultad/Facultad/Orange/www.ailab.si/orange/datasets.asp
[fail]
Error 1406 (22001): Data too long for column 'extension' at row 1
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x5e3efc]

goroutine 1 [running]:
main.scan({0x7ffd85d0b21f, 0x11}, {0x7ffd85d0b226, 0x9}, 0x4c9637?)
        /home/mario/src/nostalgia/cmd/scanner/main.go:252 +0x2fdc
main.main()
        /home/mario/src/nostalgia/cmd/scanner/main.go:342 +0x91f
exit status 2


suspect: /media/darkforce/documents/backups/yani-facultad/Facultad/Orange/

datasets.asp@Inst=on&Atts=on&Class=on&Values=on&Description=on&sort=Data+Set




curl -s 'http://localhost:3001/v1/search?q=a&type=doc&after=2010-07-31&before=2024-01-12&page=1&per_page=50' | jq .

contains=[keyword]

type:
image
doc
sheet
audio
video
zip
*any

only_directories=true / *false

date-modified=after-before

pagination= page/ per_page

{
  "directories": [
  ],
  "files": [
  ],
  "page": 1,
  "per_page": 5,
  "query": "a",
  "total_directories": 1,
  "total_files": 2
}

{
  "results": [
  ],
  "page": 1,
  "per_page": 5,
  "contains": "a",
  "total": 100
}

// file
{
  "id": 100829,
  "name": "PXL_20221225_050620590.PORTRAIT.jpg",
  "extension": "jpg",
  "path": "1/pictures/2023-03-12_pixel-7-pro",
  "date_modified": "01-01-1970",
  "size": "6.4 MB",
  "hash": "89602e0e10c7df96b624507a8c0b7af1af912373"
}

// directory
{
  "id": 1535,
  "name": "mario-facultad",
  "path": "1/documents/years/2000/mario-facultad",
  "date_modified": "06-10-1995",
  "size": "0.0 Bytes",
  "file_count": 0,
  "directory_count": 3,
  "parent_id": 1175
}