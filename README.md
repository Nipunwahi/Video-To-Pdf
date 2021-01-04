# Video-To-Pdf

 ## Installation
1. go get golang.org/x/image/bmp
2. go get github.com/jung-kurt/gofpdf
3. FFmpeg is required to run the application
4. go build video.go

## How to Use
- run the exectible with arguments <br>
<nbsp><nbsp><nbsp><nbsp>p for path , i for interval , s for size , ss for starttime , o for output file 
- only p is mandatory other arguments are optional
  
Example: This will take a frame every 20 seconds interval and make a pdf file with the outfile.pdf
```
./video -p filename.mkv -i 20 -o outfile.pdf
```
