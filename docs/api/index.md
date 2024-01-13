# api

## ping 

GET /v1/ping
> returns 'pong', health status

Responses: 200

request example

    curl -s 'http://localhost:3001/v1/ping' | jq .

response example

    {
      "message": "pong"
    }

## search

GET /v1/search?q=a&type=[image|doc|sheet|audio|video|zip|any]&after=1000-01-01&before=9999-12-31&page=1&per_page=50
> search for a resource, query, type, date after & before and paging

Responses: 200, 404

request example

    curl -s 'http://localhost:3001/v1/search?q=a&type=doc&after=2010-07-31&before=2024-01-12&page=1&per_page=50' | jq .

response example

    {
      "directories": [
        {
        "id": 2,
        "name": "docs",
        "path": "1/etc/docs/",
        "date_modified": "07-02-2023",
        "size": "6.1 kB",
        "file_count": 1,
        "directory_count": 1,
        "parent_id": 1
        }
      ],
      "files": [
        {
          "id": 2,
          "name": "Design Patterns Elements of Reusable Object-Oriented Software.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "02-08-2012",
          "size": "4.3 MB",
          "hash": "8865aeb8efaa49a1700230e2cb1dee4c157800c8"
        },
        {
          "id": 1,
          "name": "Building_Maintainable_Software_SIG.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "11-10-2016",
          "size": "6.5 MB",
          "hash": "3cf2bebbdadfe1a9fb6112c102553db0f1d7ed9b"
        }
      ],
      "page": 1,
      "per_page": 5,
      "query": "a",
      "total_directories": 1,
      "total_files": 2
    }

error codes

    404 if the query return no results

## files 

GET /v1/files/[id]

> files representation

Responses: 200, 404

request example 

    curl -s 'http://localhost:3001/v1/files/1' | jq .

response example

    {
      "id": 3,
      "name": "coreutils.pdf",
      "extension": "pdf",
      "path": "1/doc",
      "date_modified": "04-01-2023",
      "size": "1.2 MB",
      "hash": "40877fd288bc8c6118518d6c5fe565d67658d24e"
    }

error codes

    404 if there is no file with id [id]

## directories

GET /v1/directories/[id]

> directories representation

Responses: 200, 404

request example

    curl -s 'http://localhost:3001/v1/directories/3' | jq .

response example

    {
      "id": 3,
      "name": "doc",
      "path": "1/doc",
      "date_modified": "21-05-2023",
      "size": "12.0 MB",
      "file_count": 3,
      "directory_count": 0,
      "parent_id": 2,
    }

error codes

    404 if there is no directory with id [id]
 
## filescount

GET /v1/filescount

> returns the total database file count

Responses: 200

request example

    curl -s 'http://localhost:3001/v1/filescount' | jq .

response example

    {
      "count": 5
    }

## directory files

GET /v1/directories/[id]/files

> files from directory with specific id

Responses: 200, 404

request example 

    curl -s 'http://localhost:3001/v1/directories/2/files' | jq .

response example

    [
        {
          "id": 2,
          "name": "Design Patterns Elements of Reusable Object-Oriented Software.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "02-08-2012",
          "size": "4.3 MB",
          "hash": "8865aeb8efaa49a1700230e2cb1dee4c157800c8"
        },
        {
          "id": 1,
          "name": "Building_Maintainable_Software_SIG.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "11-10-2016",
          "size": "6.5 MB",
          "hash": "3cf2bebbdadfe1a9fb6112c102553db0f1d7ed9b"
        }        
    ]

error codes

    404 if there is no directory with id [id]

## directory directories

GET /v1/directories/[id]/directories

> directories from directory with specific id

Responses: 200, 404

request example 

    curl -s 'http://localhost:3001/v1/directories/2/directories' | jq .

response example

    [
        {
        "id": 3,
        "name": "doc",
        "path": "1/doc",
        "date_modified": "21-05-2023",
        "size": "12.0 MB",
        "file_count": 3,
        "directory_count": 0,
        "parent_id": 2,
        }       
    ]

error codes

    404 if there is no directory with id [id]