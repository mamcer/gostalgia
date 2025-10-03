# Nostalgia

homunculus

duplicated_GB
799.15GB

size_GB (unique files)
1792.50

## file structure

- cd  
- iso  
- year  
- music  
- picture  
- video  

example:

```
cd
    backup-01
    util-utilitarios-01
iso
    2016-03-31_eco-ciro.iso
year
    2016
music
picture
    2023-03-12_pixel-7-pro
video        
```

## process

- read-directories
- file-structure (directories, files)
- hash
- update file size
- check-existing (exists, existing-id)
- persist
- copy-files

example

    scan --tags="tag01,tag 02,tag03" --source="cd|iso|year|music|picture|video"

    config.json
    	db connection
    	scan dir
    	nostalgia dir

