## WEBSCRAPPER IN GO DEVELOPED BY OAU
This  app is using Go version 1.17 because I think for now this is the most updated release. 

### Installation
 
```
$ git clone https://github.com/omuha/oauwebscrapper.git

$ cd oauwebscrapper

$ go get
```

### Run Development Mode

- when in folder oauwebscrapper run the command bellow
    ```
    $ go run main.go
    ```
- app will exposed to this url localhost:3000/books


### Run Production with Docker

```bash
docker build -t webscrapper .
docker run -d -p 3000:3000 webscrapper
```