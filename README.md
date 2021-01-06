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

## How it works
1. It uses ffmpeg to get the frames of the video at specific intervals.
2. After getting frames, it first compares if two adjacent frames are similar or not. Similarity is found by breaking the image in some segments and getting a checksum for it.
If the images are similar the checksum would be same in most segments.
3. Using only unique frames , we store them in temp directory and then use gofpdf to make use them to make pdf.
